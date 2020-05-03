package job

import (
	"bufio"
	"fmt"
	"log"
	"nfCollector/pkg/cnf"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

const INTERVAL_PERIOD time.Duration = 1 * time.Second

const HOUR_TO_TICK int = 00
const MINUTE_TO_TICK int = 01
const SECOND_TO_TICK int = 00

// An uninteresting service.
type Job struct {
	ch        chan bool
	waitGroup *sync.WaitGroup
	timer     *time.Timer
}

// Make a new Service.
func NewJob() *Job {
	j := &Job{
		ch:        make(chan bool),
		waitGroup: &sync.WaitGroup{},
	}
	return j
}

func (j *Job) updateTimer() {
	nextTick := time.Date(time.Now().Year(), time.Now().Month(),
		time.Now().Day(), HOUR_TO_TICK, MINUTE_TO_TICK, SECOND_TO_TICK, 0, time.Local)
	if !nextTick.After(time.Now()) {
		nextTick = nextTick.Add(INTERVAL_PERIOD)
	}
	//log.Println("Next job interval is:", nextTick)
	diff := nextTick.Sub(time.Now())
	if j.timer == nil {
		j.timer = time.NewTimer(diff)
	} else {
		j.timer.Reset(diff)
	}
}

func (j *Job) Run() {
	defer j.waitGroup.Done()
	j.waitGroup.Add(1)

	log.Println("Export job started")

	j.updateTimer()

	go func() {
		defer j.waitGroup.Done()
		j.waitGroup.Add(1)
		for {
			<-j.timer.C
			// Runs Every Minute
			if time.Now().Second() == 01 {
				time.Sleep(3 * time.Second)
				conf, err := cnf.ReadConfig()
				if err != nil {
					panic(err)
				}
				if conf.Exporter.Enable == false {
					return
				}
				switch conf.Exporter.Type {
				case "influxdb":
					log.Printf("Export job configs to send metrics to InfluxDB: %v:%v (Database: %v)\n", conf.InfluxDB.Host, conf.InfluxDB.Port, conf.InfluxDB.Database)
					go j.WriteToDb(conf.Measurements.SummaryProto)
					go j.WriteToDb(conf.Measurements.SummaryProtoGeo)
				case "another":
					return
				}
			}
			j.updateTimer()
		}
	}()

}

// Stop the service by closing the service's channel.  Block until the service
// is really stopped.
func (j *Job) Stop() {
	//j.waitGroup.Wait()
	//close(s.ch)
	log.Println("Jobs Stopped!")
}

// WriteSumProto - Read File
func (j *Job) WriteToDb(meas string) {
	defer j.waitGroup.Done()
	j.waitGroup.Add(1)

	conf, err := cnf.ReadConfig()
	if err != nil {
		panic(err)
	}

	log.Printf("Job for measurement (%v) Started!", meas)
	defer log.Printf("Job for measurement (%v) Finished!", meas)

	dirname := conf.InfluxDB.TmpDir + meas + string(os.PathSeparator)
	d, err := os.Open(dirname)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer d.Close()

	// Add done extension to .working
	deleteTheDoneExtensions(dirname)

	files, err := d.Readdir(-1)
	if err != nil {
		log.Fatal("Error in opening directory", err)
	}

	log.Println("Reading directory", dirname)

	var rawMetrics []string

	// Read Directory
	for _, file := range files {
		if file.Mode().IsRegular() {
			min, _ := strconv.Atoi(strings.Split(file.Name(), "-")[3])
			if min == time.Now().Minute() {
				continue
			}
			if filepath.Ext(file.Name()) == ".metrics" {
				log.Println("Reading metrics", file.Name())
				f, err := os.OpenFile(dirname+file.Name(), os.O_RDONLY, 0600)
				if err != nil {
					log.Println("File reading error", dirname+file.Name(), err)
					continue
				}

				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					l := scanner.Text()
					if strings.Trim(l, " ") != "" {
						rawMetrics = append(rawMetrics, l)
					}
				}
				if err := scanner.Err(); err != nil {
					log.Println("File reading error", dirname+file.Name(), err)
					continue
				}

				f.Close()

				err = os.Rename(dirname+file.Name(), dirname+file.Name()+".working")
				if err != nil {
					log.Fatal("Can not rename file", dirname+file.Name(), err)
				}

			}
		}
	}

	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:               fmt.Sprintf("http://%s:%v", conf.InfluxDB.Host, conf.InfluxDB.Port),
		Username:           conf.InfluxDB.Username,
		Password:           conf.InfluxDB.Password,
		UserAgent:          "",
		Timeout:            0,
		InsecureSkipVerify: false,
		TLSConfig:          nil,
		Proxy:              nil,
	})
	if err != nil {
		log.Fatal("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database: conf.InfluxDB.Database,

		// LET InfluxDB Configuration choose the Precision
		// Precision: "u",
	})

	for _, rawMet := range rawMetrics {
		mes := strings.Split(strings.Split(rawMet, " ")[0], ",")[0]
		tags := strings.Split(strings.Split(rawMet, " ")[0], ",")[1:]
		fields := strings.Split(strings.Split(rawMet, " ")[1], ",")

		// TimeStamp
		ts := strings.Split(rawMet, " ")[2]
		msInt, _ := strconv.ParseInt(ts, 10, 64)
		timeStamp := time.Unix(0, msInt)

		tg := make(map[string]string)
		for _, t := range tags {
			tName := strings.Split(t, "=")[0]
			tVal := strings.Split(t, "=")[1]
			tg[tName] = tVal
		}

		fld := make(map[string]interface{})
		for _, f := range fields {
			fName := strings.Split(f, "=")[0]
			fVal := strings.Split(f, "=")[1]
			if "i" == fVal[len(fVal)-1:] {
				fVal = fVal[:len(fVal)-1]
				fValInt32, _ := strconv.ParseInt(fVal, 0, 32)
				fld[fName] = fValInt32
			} else if "f" == fVal[len(fVal)-1:] {
				fVal = fVal[:len(fVal)-1]
				fValInt32, _ := strconv.ParseFloat(fVal, 32)
				fld[fName] = fValInt32
			} else {
				fld[fName] = fVal
			}

		}

		pt, err := client.NewPoint(
			mes,
			tg,
			fld,
			//time.Now(),
			timeStamp,
		)

		if err != nil {
			log.Fatal("Error:", err.Error())
			continue
		}

		bp.AddPoint(pt)
	}
	err = c.Write(bp)
	if err != nil {
		log.Fatal("Error: ", err.Error())
	}

	// Add done extension to .working
	doneTheExtension(dirname)

}

// deleteTheExtensions
func deleteTheDoneExtensions(dirname string) {
	d, err := os.Open(dirname)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		log.Fatal("Error in opening directory", err)
	}

	log.Println("Deleting .done extension & rewrite suspended .working files in directory: ", dirname)

	// Read Files and rename
	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ".done" {
				err = os.Remove(dirname + file.Name())
				if err != nil {
					log.Fatal("Cant delete to .done", dirname+file.Name())
				}
				log.Printf("File %v deteled", file.Name())
			}

			if filepath.Ext(file.Name()) == ".working" {
				newFileNameArr := strings.Split(file.Name(), ".")
				err = os.Rename(dirname+file.Name(), dirname+newFileNameArr[0]+"."+newFileNameArr[1])
				if err != nil {
					log.Println("Cant rename from .working to .metrics", dirname+file.Name())
					continue
				}
				log.Printf("File %v renamed to %v", file.Name(), newFileNameArr[0]+"."+newFileNameArr[1])
			}
		}
	}
	time.Sleep(2 * time.Second)
}

// doneTheExtension
func doneTheExtension(dirname string) {
	d, err := os.Open(dirname)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		log.Fatal("Error in opening directory", err)
	}

	log.Println("Adding .done extension to .metrics in directory: ", dirname)

	// Read Files and rename
	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ".working" {
				err = os.Rename(dirname+file.Name(), dirname+file.Name()+".done")
				if err != nil {
					log.Fatal("Cant rename to .done", dirname+file.Name())
				}
				log.Printf("File %v renamed to %v", file.Name(), file.Name()+".done\n")
			}
		}
	}
}

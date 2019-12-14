package fwdr_dmp

import (
	"bytes"
	"encoding/gob"
	"github.com/landoop/tableprinter"
	"nfCollector/pkg/utl"
	"os"

	"log"
	"net"
	"sync"
)

const bufferSize int = 8960

type FlowDumper struct {
	ch					chan bool
	waitGroup 			*sync.WaitGroup
	UserFilter			Filter
}


// Make new Flow Dumper
func NewFlowDumber(fil Filter) *FlowDumper {
	fd := &FlowDumper{
		ch:        make(chan bool),
		waitGroup: &sync.WaitGroup{},
		UserFilter: fil,
	}
	return fd
}

func (fd *FlowDumper) Serve(address string, port string) {
	defer fd.waitGroup.Done()

	sAddr, err := net.ResolveUDPAddr("udp", address + ":" + port)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", sAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	if err = conn.SetReadBuffer(bufferSize); err != nil {
		log.Fatal(err)
	}
	log.Printf("Dumper is listening on %s\n", conn.LocalAddr().String())

	data := make([]byte, bufferSize)

	for {

		select {
		case <-fd.ch:
			log.Println("Disconnecting! Please wait...",)
			return
		default:
		}


		length, _, err := conn.ReadFrom(data)
		//log.Fatal("Received from UDP client :  ", string(data[:length]))

		tmpBuff := bytes.NewBuffer(data[:length])
		tmpMetrics := new([]utl.Metric)

		gobObj := gob.NewDecoder(tmpBuff)
		err = gobObj.Decode(tmpMetrics)
		if err != nil {
			log.Println("Can not decode ", err)
			continue
		}

		go fd.dumpIt(*tmpMetrics)
	}
}


func (fd *FlowDumper) dumpIt(metrics []utl.Metric) {
	defer fd.waitGroup.Done()
	fd.waitGroup.Add(1)

	filtered := fd.UserFilter.PrepareFilteredMetrics(metrics)

	DumpFiltered(filtered)
}


// Stop the service by closing the service's channel.  Block until the service
// is really stopped.
func (fd *FlowDumper) Stop() {
	fd.waitGroup.Wait()
	//close(s.ch)
	log.Println("Listener has stopped successfully!")
}




func DumpFiltered(filtered []FilteredMetric) {
	printer := tableprinter.New(os.Stdout)
	printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = true, true, true, true
	printer.CenterSeparator = "│"
	printer.ColumnSeparator = "│"
	printer.RowSeparator = "─"
	printer.Print(filtered)
}

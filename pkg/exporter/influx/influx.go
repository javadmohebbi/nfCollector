package influx

import (
	"fmt"
	"log"
	"nfCollector/pkg/cnf"
	"nfCollector/pkg/ip2loc"
	"os"
	"strconv"
	"strings"
	"time"
	//client "github.com/influxdata/influxdb/client/v2"

	"nfCollector/pkg/utl"
)


// Write - Send metrics to influx DB
func Write(metrics []utl.Metric) {

	unixNano := time.Now().UnixNano()

	// Write Protocol summary to InfluxDB
	WriteSummaryProto(metrics, unixNano)

	// Write Protocol summary GEO to InfluxDB
	WriteSummaryProtoGeo(metrics, unixNano)



}


// WriteSummaryProto - Write Protocol Summary to InfluxDB
func WriteSummaryProto (metrics []utl.Metric, unixNano int64) {
	conf, err := cnf.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	data := ""
	for _, m := range metrics {
		data += fmt.Sprintf("%v,NfVersion=%v,ExpHost=%v,ProtoName=%v,sHost=%s,sPort=%s,dHost=%s,dPort=%s Bytes=%vi,Packets=%vi,Version=\"%v\" %v\n",
			conf.Measurements.SummaryProto, m.FlowVersion, m.NFSender, m.ProtoName, m.SrcIP, m.SrcPortName, m.DstIP, m.DstPortName, m.Bytes, m.Packets, m.FlowVersion, unixNano,
		)
	}
	ToFile(data, conf.Measurements.SummaryProto)
}

// WriteSummaryProtoGeo - Write GEO Protocol Summary to InfluxDB
func WriteSummaryProtoGeo(metrics []utl.Metric, unixNano int64) {
	conf, err := cnf.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	data := ""

	var allIPs []string
	var UniqIPs []string
	var tmpGeo = make(map[string]ip2loc.IP2Loc)

	// Add All IP addresses to Slices
	for _, m := range metrics {
		allIPs = append(allIPs, m.SrcIP)
		allIPs = append(allIPs, m.DstIP)
	}

	// Remove Duplicate IPs from list
	UniqIPs = removeDuplicateIPs(allIPs)

	for _, el := range UniqIPs {
		i2l, err := ip2loc.Run(el)
		if err != nil {
			log.Println("IP2Location Error: ", err)
			continue
		}


		i2l.CountryLong 	= strings.Replace(i2l.CountryLong, ",", " ", -1)
		i2l.CountryLong 	= strings.Replace(i2l.CountryLong, "'", " ", -1)
		i2l.CountryLong 	= strings.Replace(i2l.CountryLong, " ", "_", -1)

		i2l.CountryShort 	= strings.Replace(i2l.CountryShort, " ", "-", -1)

		i2l.City 			= strings.Replace(i2l.City, ",", " ", -1)
		i2l.City 			= strings.Replace(i2l.City, "'", " ", -1)
		i2l.City 			= strings.Replace(i2l.City, " ", "-", -1)

		i2l.State 			= strings.Replace(i2l.State, ",", " ", -1)
		i2l.State 			= strings.Replace(i2l.State, "'", " ", -1)
		i2l.State 			= strings.Replace(i2l.State, " ", "-", -1)

		tmpGeo[el] = i2l
	}

	for _, m := range metrics {
		m.SrcIp2lCountryShort 	= tmpGeo[m.SrcIP].CountryShort
		m.SrcIp2lCountryLong 	= tmpGeo[m.SrcIP].CountryLong
		m.SrcIp2lState 			= tmpGeo[m.SrcIP].State
		m.SrcIp2lCity 			= tmpGeo[m.SrcIP].City
		m.SrcIp2lLat 			= tmpGeo[m.SrcIP].Lat
		m.SrcIp2lLong 			= tmpGeo[m.SrcIP].Long

		m.DstIp2lCountryShort 	= tmpGeo[m.DstIP].CountryShort
		m.DstIp2lCountryLong 	= tmpGeo[m.DstIP].CountryLong
		m.DstIp2lState 			= tmpGeo[m.DstIP].State
		m.DstIp2lCity 			= tmpGeo[m.DstIP].City
		m.DstIp2lLat 			= tmpGeo[m.DstIP].Lat
		m.DstIp2lLong 			= tmpGeo[m.DstIP].Long

		data += fmt.Sprintf("%v,NfVersion=%v,ExpHost=%v,ProtoName=%v,sHost=%s,sPort=%s,dHost=%s,dPort=%s,sCouSh=%v,sCouLo=%v,sSta=%v,sCit=%v,dCouSh=%v,dCouLo=%v,dSta=%v,dCit=%v Bytes=%vi,Packets=%vi,Version=\"%v\",",
			conf.Measurements.SummaryProtoGeo, m.FlowVersion, m.NFSender, m.ProtoName, m.SrcIP, m.SrcPortName, m.DstIP, m.DstPortName,
			m.SrcIp2lCountryShort, m.SrcIp2lCountryLong, m.SrcIp2lState, m.SrcIp2lCity,
			m.DstIp2lCountryShort, m.DstIp2lCountryLong, m.DstIp2lState, m.DstIp2lCity,
			m.Bytes, m.Packets, m.FlowVersion,

		)

		dLat, _ := strconv.ParseFloat(m.DstIp2lLat, 64)
		dLon, _ := strconv.ParseFloat(m.DstIp2lLong, 64)

		sLat, _ := strconv.ParseFloat(m.SrcIp2lLat, 64)
		sLon, _ := strconv.ParseFloat(m.SrcIp2lLong, 64)

		data += fmt.Sprintf("ddLat=%ff,ddLon=%ff,ssLat=%ff,ssLon=%ff %v\n",
			dLat, dLon,
			sLat, sLon, unixNano)

	}

	ToFile(data, conf.Measurements.SummaryProtoGeo)
}



// removeDuplicateIPs - Remove Duplicate entries from IP address slice
func removeDuplicateIPs(stringSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// ToFile - Export Metrics to File
func ToFile(data string, measurementName string){
	conf, err := cnf.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	currentTime := time.Now()
	dtSpan := currentTime.Format("20060102-15-04")
	fileName := fmt.Sprintf("nfc-%v-(%v).metrics", dtSpan, measurementName)

	dir := conf.InfluxDB.TmpDir + measurementName + string(os.PathSeparator)

	_ = os.MkdirAll(conf.InfluxDB.TmpDir + measurementName, os.ModePerm)

	f, err := os.OpenFile(dir + fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err = f.WriteString(data); err != nil {
		log.Fatal(err)
	}

	return
}
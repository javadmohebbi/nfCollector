package exporter

import (
	"nfCollector/pkg/cnf"
	"nfCollector/pkg/exporter/influx"
	"nfCollector/pkg/utl"
)

// Write to the selected DB
func Write(m []utl.Metric){
	conf, err := cnf.ReadConfig()
	if err != nil {
		panic(err)
	}

	if conf.Exporter.Enable == false {
		return
	}
	switch conf.Exporter.Type {
	case "influxdb":
		ExportToInflux(m)
	case "another":
		return
	}
}

// ExportToInflux
func ExportToInflux(m []utl.Metric) {
	//if influx.InfluxDBConn == nil {
	//	influx.InfluxDBConn = influx.Connect()
	//}
	influx.Write(m)
}


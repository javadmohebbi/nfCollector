package main

import (
	"flag"
	"log"
	fwdr_dmp "nfCollector/pkg/fwdr-dmp"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	argAddress := flag.String("addr", "", "Listen IP address - Default 127.0.0.1")
	argPort := flag.String("port", "", "Listen port - Default 7161")

	argFilterVersion 		:= flag.String("flt-nf-ver", "*", "Filter netflow version. eg: 1, 5, 6, 7, 9, 10 (for IPFIX)")
	argFilterExporter 		:= flag.String("flt-nf-exp", "*", "Filter netflow Exporter. IP address of exporter device. eg: 192.168.1.1, 192.168.1.*")
	argFilterSrcIP 			:= flag.String("flt-src-ip", "*", "Filter Source IP. eg: 192.168.1.1, 192.168.1.*")
	argFilterSrcPort 		:= flag.String("flt-src-port", "*", "Filter Source Port. eg: 80, 433, 100-250")
	argFilterDstIP 			:= flag.String("flt-dst-ip", "*", "Filter Destination IP. eg: 192.168.1.1, 192.168.1.*")
	argFilterDstPort 		:= flag.String("flt-dst-port", "*", "Filter Destination Port. eg: 80, 433, 100-250")
	argFilterProtocol 		:= flag.String("flt-proto", "*", "Filter Protocol. eg: tcp, udp, icmp")

	// Parse command line options
	flag.Parse()

	filter := fwdr_dmp.NewFilter(*argFilterVersion, *argFilterExporter, *argFilterSrcIP, *argFilterSrcPort, *argFilterDstIP, *argFilterDstPort, *argFilterProtocol)

	fldu := fwdr_dmp.NewFlowDumber(*filter)

	addr := "127.0.0.1"
	port := "7161"

	if *argAddress != ""{
		addr = *argAddress
	}

	if *argPort != ""{
		port = *argPort
	}

	go fldu.Serve(addr, port)

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	// Stop the service gracefully.
	fldu.Stop()
}
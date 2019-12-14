package flgs

import (
	"flag"
	"fmt"

	"log"
	"net"
	"nfCollector/pkg/cnf"
	"nfCollector/pkg/job"
	"nfCollector/pkg/lstn"
	"nfCollector/pkg/service"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	version        = "0.0.1"
	appName        = "Netflow Collector"
	websiteAddress = "http://nfc.mjmohebbi.com"
)

// InitFlags - Initialize input options
func InitFlags() {

	// Define Command line options
	argVer := flag.Bool("v", false, "Print Version & exit.")
	argAddress := flag.String("addr", "", "Listen IP address")
	argPort := flag.String("port", "", "Listen port")
	argDump := flag.String("dump", "", "It will Print flow record if the value is 'true' and 'false' for nothing")
	argDebug := flag.String("debug", "", "It will Print debug info if the value is 'true' and 'false' for nothing")

	// Parse command line options
	flag.Parse()

	debugInfo := false
	// Debug Info
	if *argDebug == "true" {
		debugInfo = true
	}

	// Print version
	if *argVer {
		PrintVersion()
	}

	// Listen & Extract
	con, isDump := ListenOn(argAddress, argPort, argDump)

	jb := job.NewJob()
	go jb.Run()

	conf, err := cnf.ReadConfig()
	if err != nil {
		log.Fatal("Can not read config", err)
	}

	svc := service.NewService(*jb)
	go svc.Serve(con, isDump, debugInfo, conf.Server.Forwarder)

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	svc.Stop()

}

// PrintVersion - Print App Version, Info & Exit
func PrintVersion() {
	fmt.Printf("%s | %s | %s", appName, version, websiteAddress)
	os.Exit(0)
}

// ListenOn - Call Listen & Extract UDP packets
func ListenOn(addr *string, port *string, argDump *string) (*net.UDPConn, bool) {
	listenIP := ""
	listenPort := ""
	isDump := false
	config, err := cnf.ReadConfig()
	if err != nil {
		panic(err)
	}
	if *addr == "" && err == nil {
		// Read from config
		listenIP = config.Server.Address
	} else {
		listenIP = *addr
	}

	if *port == "" && err == nil {
		listenPort = strconv.Itoa(config.Server.Port)
	} else {
		listenPort = *port
	}

	if *argDump == "" && err == nil {
		isDump = config.Server.Dump
	} else {
		isDump, _ = strconv.ParseBool(*argDump)
	}

	if listenIP == "" || listenIP == "0.0.0.0" {
		listenIP = ""
	}

	ln, err := lstn.Listen(listenIP, listenPort)

	if err != nil {
		log.Println(err)
		log.Println("You should kill running 'nfc' instance to start it again!")
		os.Exit(1)
	}

	return ln, isDump
}

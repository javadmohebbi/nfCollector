package service

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"nfCollector/pkg/cnf"
	"nfCollector/pkg/exporter"
	"nfCollector/pkg/job"
	"nfCollector/pkg/nfipfix"
	"nfCollector/pkg/nfv1"
	"nfCollector/pkg/nfv5"
	"nfCollector/pkg/nfv6"
	"nfCollector/pkg/nfv7"
	"nfCollector/pkg/nfv9"
	"nfCollector/pkg/utl"
	"strconv"
	"strings"
	"sync"

	"github.com/tehmaze/netflow"
	"github.com/tehmaze/netflow/ipfix"
	"github.com/tehmaze/netflow/netflow1"
	"github.com/tehmaze/netflow/netflow5"
	"github.com/tehmaze/netflow/netflow6"
	"github.com/tehmaze/netflow/netflow7"
	"github.com/tehmaze/netflow/netflow9"
	"github.com/tehmaze/netflow/session"
)

const bufferSize int = 8960

// Service - An uninteresting service.
type Service struct {
	ch        chan bool
	waitGroup *sync.WaitGroup
	job       *job.Job
}

// NewService - Make a new Service.
func NewService(j job.Job) *Service {
	s := &Service{
		ch:        make(chan bool),
		waitGroup: &sync.WaitGroup{},
		job:       &j,
	}
	return s
}

// Serve - Accept connections and spawn a goroutine to serve each one.  Stop listening
// if anything is received on the service's channel.
func (s *Service) Serve(conn *net.UDPConn, dump bool, debug bool, forward bool) {
	defer s.waitGroup.Done()
	defer conn.Close()

	data := make([]byte, bufferSize)
	decoders := make(map[string]*netflow.Decoder)

	for {

		select {
		case <-s.ch:
			log.Println("Disconnecting! Please wait...")
			return
		default:
		}

		length, remote, err := conn.ReadFrom(data)
		if err != nil {
			log.Println(err)
			continue
		}
		d, found := decoders[remote.String()]
		if !found {
			s := session.New()
			d = netflow.NewDecoder(s)
			decoders[remote.String()] = d
		}
		m, err := d.Read(bytes.NewBuffer(data[:length]))
		if err != nil {
			if !dump {
				log.Println("decoder error:", err)
			}
			continue
		}
		if debug {
			log.Printf("received %d bytes from %s\n", length, remote)
		}

		go s.checkPacket(m, remote, dump, data, forward)

	}

}

// Stop the service by closing the service's channel.  Block until the service
// is really stopped.
func (s *Service) Stop() {
	s.waitGroup.Wait()
	s.job.Stop()
	//close(s.ch)
	log.Println("Listener has stopped successfully!")
}

func (s *Service) checkPacket(m interface{}, remote net.Addr, dump bool, data []byte, isForward bool) {

	defer s.waitGroup.Done()
	s.waitGroup.Add(1)

	var metrics []utl.Metric

	switch p := m.(type) {
	case *netflow1.Packet:
		metrics = nfv1.Prepare(remote.String(), p)

	case *netflow5.Packet:
		metrics = nfv5.Prepare(remote.String(), p)

	case *netflow6.Packet:
		metrics = nfv6.Prepare(remote.String(), p)

	case *netflow7.Packet:
		metrics = nfv7.Prepare(remote.String(), p)

	case *netflow9.Packet:
		metrics = nfv9.Prepare(remote.String(), p)

	case *ipfix.Message:
		metrics = nfipfix.Prepare(remote.String(), p)
	}

	// Dump tp Standard Output
	if dump {
		utl.Dump(metrics)
	} else {
		// Export
		exporter.Write(metrics)
	}

	// Forward
	if isForward {
		conf, err := cnf.ReadConfig()
		if err != nil {
			log.Fatal("Can not read config", err)
		}
		hosts := strings.Split(conf.Server.ForwarderHost, ";")
		port := strconv.Itoa(conf.Server.ForwarderPort)
		for _, host := range hosts {
			host = strings.TrimSpace(host)
			// Forward UDP Packet
			go s.forwardTo(metrics, host, port)
		}
	}
}

func (s *Service) forwardTo(metrics []utl.Metric, host string, port string) {

	defer s.waitGroup.Done()
	s.waitGroup.Add(1)

	//log.Println(host + ":" + port)
	fw, err := net.ResolveUDPAddr("udp", host+":"+port)
	if err != nil {
		log.Println("Can not forward logs due to error", err)
		return
	}

	fwCon, err := net.DialUDP("udp", nil, fw)
	if err != nil {
		log.Println("Can not forward logs due to error", err)
		return
	}

	defer fwCon.Close()

	binBuf := new(bytes.Buffer)
	gobObj := gob.NewEncoder(binBuf)
	err = gobObj.Encode(metrics)
	if err != nil {
		log.Println("Err: Can not Encode metrics", err)
	}

	_, err = fwCon.Write(binBuf.Bytes())
	if err != nil {
		log.Println("Can not forward logs due to error", err)
		return
	}

}

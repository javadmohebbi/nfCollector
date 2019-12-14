package lstn

import (
	"log"
	"net"
	"os"
)

const bufferSize int = 8960




// Listen - Listen Server on proto, address:port
func Listen(address string, port string) (*net.UDPConn, error) {
	// Resolve Address
	sAddr, err := net.ResolveUDPAddr("udp", address + ":" + port)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Listen
	sConn, err := net.ListenUDP("udp", sAddr)
	if err != nil {
		return nil, err
	}
	if err = sConn.SetReadBuffer(bufferSize); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Printf("Server is listening on %s\n", sConn.LocalAddr().String())
	return sConn, nil
}
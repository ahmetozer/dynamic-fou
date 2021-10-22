package share

import (
	"fmt"
	"net"
)

func UDPPortCheck(port int) error {

	tempPort, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("::"),
	})

	if err != nil {
		return fmt.Errorf("can't listen on port %v, %s", port, err)
	}

	err = tempPort.Close()
	if err != nil {
		return fmt.Errorf("couldn't stop listening on port %v, %s", port, err)
	}
	return nil
}

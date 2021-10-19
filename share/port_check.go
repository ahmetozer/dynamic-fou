package share

import (
	"fmt"
	"net"
	"strconv"
)

func UDPPortCheck(port uint16) error {
	tempPort, err := net.Listen("udp", ":"+strconv.FormatUint(uint64(port), 10))

	if err != nil {
		return fmt.Errorf("can't listen on port %q: %s", port, err)
	}

	err = tempPort.Close()
	if err != nil {
		return fmt.Errorf("couldn't stop listening on port %q: %s", port, err)
	}
	return nil
}

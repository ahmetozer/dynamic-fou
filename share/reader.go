package share

import (
	"bufio"
	"net"
)

func NewReader(conn *net.Conn, i int, pC chan []byte, errC chan error) {
	for {
		p := make([]byte, i)
		_, err := bufio.NewReader(*conn).Read(p)
		if err == nil {
			pC <- p
		}
		errC <- err

	}
}

package client

import (
	"fmt"
	"net"
	"time"

	"github.com/ahmetozer/dynamic-fou/share"
)

// Ask server current IP and Port
func (a SvConfig) Whoami(conn *net.Conn) (string, string, error) {

	fmt.Fprintf(*conn, "mode=whoami\nclient=%v\notk=%v", ClientName, share.NewOTK(a.ClientKey))

	var err error
	p := make([]byte, 2048)

	pC := make(chan []byte, 2048)
	errC := make(chan error)

	go share.NewReader(conn, 2048, pC, errC)

	select {
	case i := <-pC:
		p = i
	case i := <-errC:
		err = i
	case <-time.After(3000 * time.Millisecond):
		err = fmt.Errorf("respond read time out")
		break
	}

	var IP, PORT string

	if err == nil {
		IP = share.IniVal(string(p), "IP")
		PORT = share.IniVal(string(p), "PORT")
		if IP == "" || PORT == "" {
			return "", "", fmt.Errorf("server respond is corrupted ip:'%v' port:'%v'", IP, PORT)
		}
	} else {
		return "", "", err
	}
	return IP, PORT, nil
}

package client

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/ahmetozer/dynamic-fou/share"
)

type Whoami struct {
	IP              string
	PORT            string
	REMOTE_FOU_PORT int
}

// Ask server current IP and Port
func (a SvConfig) Whoami(conn *net.Conn) (Whoami, error) {

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

	var IP, PORT, FOU_PORT string

	if err == nil {
		IP = share.IniVal(string(p), "IP")
		PORT = share.IniVal(string(p), "PORT")
		if IP == "" || PORT == "" {
			return Whoami{}, fmt.Errorf("server respond is corrupted ip:'%v' port:'%v'", IP, PORT)
		}
	} else {
		return Whoami{}, err
	}

	FOU_PORT = share.IniVal(string(p), "FOU_PORT")
	if FOU_PORT == "" {
		return Whoami{}, fmt.Errorf("fou port respond is empty ip:'%v' port:'%v'", IP, PORT)
	}
	FOU_PORT_INT, err := strconv.Atoi(FOU_PORT)
	if err != nil {
		return Whoami{}, fmt.Errorf("fou port respond is empty ip:'%v' port:'%v' %v", IP, PORT, err)
	}
	return Whoami{
		IP:              IP,
		PORT:            PORT,
		REMOTE_FOU_PORT: FOU_PORT_INT,
	}, nil
}

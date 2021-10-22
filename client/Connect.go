package client

import (
	"fmt"
	"net"
	"time"

	"github.com/ahmetozer/dynamic-fou/share"
)

func (a SvConfig) Connect(conn *net.Conn) error {

	fmt.Fprintf(*conn, "mode=connect\nclient=%v\notk=%v", ClientName, share.NewOTK(a.ClientKey))

	var status string
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
	case <-time.After(5000 * time.Millisecond):
		err = fmt.Errorf("respond read time out")
		break
	}

	if err == nil {
		status = share.IniVal(string(p), "status")
		if status != "ok" && status != "reconnect" {
			return fmt.Errorf("%v", share.IniVal(string(p), "err"))
		}
	} else {
		return err
	}
	return nil
}

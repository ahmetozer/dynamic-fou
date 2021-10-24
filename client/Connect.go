package client

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/ahmetozer/dynamic-fou/share"
)

var (
	fouPortInt map[int]int
)

func (a SvConfig) Connect(conn *net.Conn, clientId int, whoami Whoami) error {

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

	newPort := fmt.Sprintf("%v", (*conn).LocalAddr())
	newPort = newPort[strings.LastIndex(newPort, ":")+1:]
	tempPort, err := strconv.Atoi(newPort)
	if err != nil {
		return fmt.Errorf("atoi: %v", err)
	}

	if fouPortInt[clientId] != 0 {
		err = share.FouDel(fouPortInt[clientId])
		if err != nil {
			return fmt.Errorf("fouDel: %v %v", err, tempPort)
		}
	}

	if share.IsInterfacesExist(clientId) {
		err = share.InterfaceDel(clientId)
		if err != nil {
			return fmt.Errorf("interfaceDel: %v", err)
		}
	}

	connOldLocal := (*conn).LocalAddr().String()

	err = (*conn).Close()
	if err != nil {
		return fmt.Errorf("connClose: %v", err)
	}

	laddr, err := net.ResolveUDPAddr("udp", connOldLocal)
	if err != nil {
		return fmt.Errorf("laddr: %v %v", err, connOldLocal)
	}
	raddr := net.UDPAddr{IP: net.ParseIP(a.RemoteAddr), Port: whoami.REMOTE_FOU_PORT}
	tempConn, err := net.DialUDP("udp", laddr, &raddr)
	if err != nil {
		return fmt.Errorf("tempConn: %v %v", err, tempPort)
	}

	for ty := 0; ty < 5; ty++ {
		fmt.Fprintf(tempConn, "mode=connect\nclient=%v\notk=%v", ClientName, share.NewOTK(a.ClientKey))
		time.Sleep(time.Millisecond * 500)
	}

	err = tempConn.Close()
	if err != nil {
		return fmt.Errorf("tempConnClose: %v", err)
	}

	err = share.FouAdd(tempPort)
	if err != nil {
		return fmt.Errorf("fouAdd: %v %v", err, tempPort)
	}
	fouPortInt[clientId] = tempPort

	err = share.InterfaceAdd(clientId, -1, a.RemoteAddr, int(whoami.REMOTE_FOU_PORT), a.MTU)

	if err != nil {
		return fmt.Errorf("interfaceAdd: %v", err)
	}
	return nil
}

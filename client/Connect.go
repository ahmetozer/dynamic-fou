package client

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/ahmetozer/dynamic-fou/share"
	"github.com/vishvananda/netlink"
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

	// Read remote tunnel source port
	serverSourcePort, err := strconv.Atoi(share.IniVal(string(p), "source_port"))
	if err != nil {
		return fmt.Errorf("atoi: %v", err)
	}

	// Read remote tunnel source port
	serverv6LocalAddr := share.IniVal(string(p), "v6_localAddr")

	newPort := fmt.Sprintf("%v", (*conn).LocalAddr())
	newPort = newPort[strings.LastIndex(newPort, ":")+1:]
	clientFouListenPort, err := strconv.Atoi(newPort)
	if err != nil {
		return fmt.Errorf("atoi: %v", err)
	}

	if fouPortInt[clientId] != 0 {
		err = share.FouDel(fouPortInt[clientId])
		if err != nil {
			return fmt.Errorf("fouDel: %v %v", err, clientFouListenPort)
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

	// Send few packets to Fou DST port
	raddr := net.UDPAddr{IP: net.ParseIP(a.RemoteAddr), Port: whoami.REMOTE_FOU_PORT}
	tempConn, err := net.DialUDP("udp", laddr, &raddr)
	if err != nil {
		return fmt.Errorf("tempConn: %v %v", err, clientFouListenPort)
	}

	for ty := 0; ty < 3; ty++ {
		fmt.Fprintf(tempConn, "mode=connect\nclient=%v\notk=%v", ClientName, share.NewOTK(a.ClientKey))
		time.Sleep(time.Millisecond * 500)
	}

	err = tempConn.Close()
	if err != nil {
		return fmt.Errorf("tempConnClose: %v", err)
	}

	// Send few packets to Fou SRC port
	// Open client dst port for server source port
	raddr = net.UDPAddr{IP: net.ParseIP(a.RemoteAddr), Port: serverSourcePort}
	tempConn2, err := net.DialUDP("udp", laddr, &raddr)
	if err != nil {
		return fmt.Errorf("tempConn: %v %v", err, clientFouListenPort)
	}

	for ty := 0; ty < 3; ty++ {
		fmt.Fprintf(tempConn2, "mode=connect\nclient=%v\notk=%v", ClientName, share.NewOTK(a.ClientKey))
		time.Sleep(time.Millisecond * 500)
	}

	err = tempConn2.Close()

	if err != nil {
		return fmt.Errorf("tempConnClose: %v", err)
	}

	err = share.FouAdd(clientFouListenPort)
	if err != nil {
		return fmt.Errorf("fouAdd: %v %v", err, clientFouListenPort)
	}

	// Store current fou port, that will be used in the next for removing existing port
	// or clear changes on exit
	fouPortInt[clientId] = clientFouListenPort

	err = share.InterfaceAdd(clientId, clientFouListenPort, a.RemoteAddr, int(whoami.REMOTE_FOU_PORT), a.MTU)
	if err != nil {
		return fmt.Errorf("interfaceAdd: %v", err)
	}
	Interface, err := netlink.LinkByName(share.InterfaceName(clientId))
	if err != nil {
		return fmt.Errorf("interfaceSelect: %v", err)
	}
	err = netlink.LinkSetUp(Interface)
	if err != nil {
		return fmt.Errorf("interfaceUp: %v", err)
	}

	addr, err := netlink.ParseAddr(serverv6LocalAddr)
	if err != nil {
		return fmt.Errorf("parseAddr: %v", err)
	}

	err = netlink.RouteAdd(&netlink.Route{LinkIndex: Interface.Attrs().Index, Dst: addr.IPNet})
	if err != nil {
		return fmt.Errorf("addrAdd: %v", err)
	}

	return nil
}

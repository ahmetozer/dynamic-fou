package share

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/vishvananda/netlink"
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

func CheckKernelFouCapability() error {

	var err error
	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
		return fmt.Errorf("net dial: %v", err)
	}

	newPort := fmt.Sprintf("%v", conn.LocalAddr())
	newPort = newPort[strings.LastIndex(newPort, ":")+1:]
	tempPort, err := strconv.Atoi(newPort)
	if err != nil {
		return fmt.Errorf("atoi: %v", err)
	}
	err = conn.Close()
	if err != nil {
		return fmt.Errorf("conn close: %v", err)
	}

	new_fou := &netlink.Fou{
		Family:    netlink.FAMILY_V4,
		Protocol:  4,
		Port:      tempPort,
		EncapType: netlink.FOU_ENCAP_DIRECT,
	}

	err = netlink.FouAdd(*new_fou)
	if err != nil {
		return fmt.Errorf("fou add: %v", err)
	}

	newtun := netlink.Iptun{
		LinkAttrs:  netlink.LinkAttrs{Name: "fouTest", MTU: 1480},
		PMtuDisc:   1,
		Local:      net.IPv4(127, 0, 0, 1),
		Remote:     net.IPv4(127, 0, 0, 1),
		EncapSport: uint16(tempPort),
		EncapDport: 5000,
		EncapType:  netlink.FOU_ENCAP_DIRECT,
	}

	err = netlink.LinkAdd(&newtun)
	if err != nil {
		err2 := netlink.FouDel(*new_fou)
		if err2 != nil {
			return fmt.Errorf("fou del: %v", err2)
		}
		return fmt.Errorf("link add: %v", err)
	}

	err = netlink.LinkDel(&newtun)
	if err != nil {
		err2 := netlink.FouDel(*new_fou)
		if err2 != nil {
			return fmt.Errorf("fou del: %v", err2)
		}
		return fmt.Errorf("link del: %v", err)
	}
	err = netlink.FouDel(*new_fou)
	if err != nil {
		return fmt.Errorf("fou del: %v", err)
	}

	return nil
}

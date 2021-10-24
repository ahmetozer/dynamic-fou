package share

import (
	"fmt"
	"net"
	"reflect"

	"github.com/vishvananda/netlink"
	"go.uber.org/zap"
)

const InterfacePrefix = "dyn"

func FouAdd(p int) error {
	new_fou := &netlink.Fou{
		Family:    netlink.FAMILY_V4,
		Protocol:  4,
		Port:      p,
		EncapType: netlink.FOU_ENCAP_DIRECT,
	}

	return netlink.FouAdd(*new_fou)
}

func FouDel(p int) error {
	new_fou := &netlink.Fou{
		Family:    netlink.FAMILY_V4,
		Protocol:  4,
		Port:      p,
		EncapType: netlink.FOU_ENCAP_DIRECT,
	}

	return netlink.FouDel(*new_fou)
}

func InterfaceAdd(id int, sourcePort int, remote string, destinationPort int, MTU int) error {
	destinationAddress := net.ParseIP(remote)
	newtun := netlink.Iptun{}
	if sourcePort == -1 {
		newtun = netlink.Iptun{
			LinkAttrs:  netlink.LinkAttrs{Name: fmt.Sprintf("%v%v", InterfacePrefix, id), MTU: MTU},
			PMtuDisc:   1,
			Remote:     destinationAddress,
			EncapDport: uint16(destinationPort),
			EncapType:  netlink.FOU_ENCAP_DIRECT,
		}
	} else {
		newtun = netlink.Iptun{
			LinkAttrs:  netlink.LinkAttrs{Name: fmt.Sprintf("%v%v", InterfacePrefix, id), MTU: MTU},
			PMtuDisc:   1,
			Remote:     destinationAddress,
			EncapSport: uint16(sourcePort),
			EncapDport: uint16(destinationPort),
			EncapType:  netlink.FOU_ENCAP_DIRECT,
		}
	}

	Logger.Debug("new tunnel", zap.String("tun", toString(newtun)))
	return netlink.LinkAdd(&newtun)
}

func InterfaceDel(id int) error {
	newtun2 := netlink.Iptun{
		LinkAttrs: netlink.LinkAttrs{Name: fmt.Sprintf("%v%v", InterfacePrefix, id)},
	}

	return netlink.LinkDel(&newtun2)
}

func toString(k interface{}) string {
	v := reflect.ValueOf(k)
	typeOfS := v.Type()

	var t string

	l := v.NumField() - 1
	for i := 0; i < l+1; i++ {
		if i < l {
			t += fmt.Sprintf("'%s':'%v',", typeOfS.Field(i).Name, v.Field(i).Interface())
		} else {
			t += fmt.Sprintf("'%s':'%v'", typeOfS.Field(i).Name, v.Field(i).Interface())
		}

	}
	return t
}

func AddrAddArr(i string, k ...string) error {
	Interface, err := netlink.LinkByName(i)
	if err != nil {
		return err
	}
	addrErrMap := make(map[string]error)

	for i := 0; i < len(k); i++ {
		addr, err1 := netlink.ParseAddr(k[i])
		if err1 != nil {
			addrErrMap[k[i]] = err1
		} else {
			err2 := netlink.AddrAdd(Interface, addr)
			if err2 != nil {
				addrErrMap[k[i]] = err2
			}
		}

	}

	extendedError := ""
	for i, e := range addrErrMap {
		extendedError += fmt.Sprintf("'%v':'%v'", i, e)
	}

	if extendedError != "" {
		return fmt.Errorf("%v", extendedError)
	}

	return nil
}

func IsInterfacesExist(iface int) bool {
	interfaces, _ := net.Interfaces()
	for _, i := range interfaces {
		if i.Name == fmt.Sprintf("%v%v", InterfacePrefix, iface) {
			return true
		}

	}
	return false
}

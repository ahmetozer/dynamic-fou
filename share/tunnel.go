package share

import "github.com/vishvananda/netlink"

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

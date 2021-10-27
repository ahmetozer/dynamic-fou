package server

import (
	"fmt"
	"net"

	"github.com/ahmetozer/dynamic-fou/share"
	"github.com/vishvananda/netlink"
	"go.uber.org/zap"
)

func Connect(conn *net.UDPConn, remote *net.UDPAddr, client ClientConfig) {
	err := share.UDPPortCheck(remote.Port)
	status := "ok"
	if err != nil {
		status = "err"
		share.Logger.Debug("connect.PortCheck", zap.String("remote", remote.String()), zap.Error(err))
		_, err = conn.WriteTo([]byte(fmt.Sprintf("status=%v\nerr=%v\n", status, err)), remote)
		if err != nil {
			share.Logger.Error("connect", zap.String("remote", remote.String()), zap.Error(err))
		}
		status = "err"
		return
	}

	if CurrentClientList[client.ClientName].IP != "" {
		share.Logger.Debug("connect.InterfaceDel", zap.String("clientName", client.ClientName), zap.Error(share.InterfaceDel(CurrentClientIdList[client.ClientName])))
	}

	CurrentClientList[client.ClientName] = CurrentClient{
		IP:   remote.IP.String(),
		PORT: remote.Port,
	}

	sourcePort := share.PortGenerate()
	err = share.InterfaceAdd(CurrentClientIdList[client.ClientName], sourcePort, remote.IP.String(), remote.Port, client.MTU)

	if err != nil {
		share.Logger.Error("connect.InterfaceAdd", zap.String("clientName", client.ClientName), zap.Error(err))
		status = err.Error()
	}

	Interface, err := netlink.LinkByName(share.InterfaceName(CurrentClientIdList[client.ClientName]))
	if err != nil {
		share.Logger.Error("connect.InterfaceSelect", zap.String("clientName", client.ClientName), zap.Error(err))
		status = err.Error()
	}

	err = netlink.LinkSetUp(Interface)
	if err != nil {
		share.Logger.Error("connect.InterfaceUp", zap.String("clientName", client.ClientName), zap.Error(err))
		status = err.Error()
	}

	_, err = conn.WriteTo([]byte(fmt.Sprintf("status=%v\nsource_port=%v\n", status, sourcePort)), remote)
	if err != nil {
		share.Logger.Error("connect", zap.String("remote", remote.String()), zap.Error(err))
	}

}

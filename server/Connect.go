package server

import (
	"fmt"
	"net"

	"github.com/ahmetozer/dynamic-fou/share"
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
		fmt.Printf("\n\nremove client %v:%v\n\n", CurrentClientList[client.ClientName].IP, CurrentClientList[client.ClientName].PORT)
		share.Logger.Debug("connect.InterfaceDel", zap.String("clientName", client.ClientName), zap.Error(share.InterfaceDel(CurrentClientIdList[client.ClientName])))
	}

	CurrentClientList[client.ClientName] = CurrentClient{
		IP:   remote.IP.String(),
		PORT: remote.Port,
	}
	fmt.Printf("\n\n%v\n\n", CurrentClientList[client.ClientName].toString())

	//err = share.InterfaceAdd(CurrentClientIdList[client.ClientName], 999, remote.IP.String(), remote.Port, client.MTU)
	err = share.InterfaceAdd(0, 999, remote.IP.String(), 800, 1500)
	if err != nil {
		share.Logger.Error("connect.InterfaceAdd", zap.String("clientName", client.ClientName), zap.Error(err))
		status = err.Error()
	}

	_, err = conn.WriteTo([]byte(fmt.Sprintf("status=%v\n", status)), remote)
	if err != nil {
		share.Logger.Error("connect", zap.String("remote", remote.String()), zap.Error(err))
	}

}

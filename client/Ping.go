package client

import (
	"fmt"
	"net"
	"strings"
	"syscall"
	"time"

	"github.com/ahmetozer/dynamic-fou/share"
)

func (a *SvConfig) Ping(clientId int) error {
	servAddr := fmt.Sprintf("[%v%%%v]:%v", strings.Split(a.RemoteV6LocalAddr, "/")[0], share.InterfaceName(clientId), a.RemotePort)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		return fmt.Errorf("dial error: %v", err)
	}

	conn, err := net.DialTCP("tcp6", nil, tcpAddr)
	if err != nil {
		return fmt.Errorf("dial error: %v", err)
	}
	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(time.Second)
	sockFile, sockErr := conn.File()
	if sockErr == nil {
		tcpFd := int(sockFile.Fd())
		err := syscall.SetsockoptInt(tcpFd, syscall.IPPROTO_TCP, syscall.TCP_KEEPCNT, 3)
		if err != nil {
			return fmt.Errorf("setting keepalive probe count: %v", err)
		}
		err = syscall.SetsockoptInt(tcpFd, syscall.IPPROTO_TCP, syscall.TCP_KEEPINTVL, 5)
		if err != nil {
			return fmt.Errorf("setting keepalive retry interval: %v", err)
		}
		sockFile.Close()
	} else {
		return fmt.Errorf("setting socket keepalive: %v", err)
	}

	reply := make([]byte, 1024)

	_, err = conn.Read(reply)
	if err != nil {
		return fmt.Errorf("read error: %v", err)
	}

	return conn.Close()
}

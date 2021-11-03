package client

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/ahmetozer/dynamic-fou/share"
	"go.uber.org/zap"
)

var (
	ClientName string
	ScriptFile string
)

func Start() {
	fouPortInt = make(map[int]int)

	configFile := os.Getenv("SERVER_LIST")
	ScriptFile = os.Getenv("SCRIPT_FILE")

	if configFile == "" {
		configFile = "/etc/dynamic-fou.client.json"
	}
	share.Logger.Debug("opening config file", zap.String("config-file", configFile))
	err := share.CheckFolder(filepath.Dir(configFile))
	if err != nil {
		share.Logger.Fatal("log path is not oppened", zap.String("path", filepath.Dir(configFile)), zap.String("err", err.Error()))
	}

	clients, err := Parse(configFile)
	if err != nil {
		share.Logger.Fatal("config file cannot parsed", zap.String("file", configFile), zap.String("err", err.Error()))
	}

	share.Logger.Debug("config file parsed", zap.String("clientName", clients.ClientName), zap.String("server-count", fmt.Sprint((len(clients.Servers)))))
	ClientName = clients.ClientName
	var wg sync.WaitGroup
	for i := 0; i < len(clients.Servers); i++ {
		wg.Add(1)
		share.Logger.Debug(fmt.Sprintf("client %v", i+1), clients.Servers[i].toZap()...)
		go clients.Servers[i].clientInit(i + 1)
	}

	share.Logger.Debug("client inits are done.")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		share.Logger.Debug("shutdown: Deleting Fou endpoints and tunnels")
		os.Exit(Shutdown(clients))
	}()

	wg.Wait()
}

func (a SvConfig) clientInit(clientId int) {
	defer func() {
		if r := recover(); r != nil {
			share.Logger.Error("panic", zap.Int("clientId", clientId), zap.String("recover", fmt.Sprintf("%v", r)))
		}
	}()
	k := true
	var whoami Whoami
INITLOOP:
	for k {
		remote := fmt.Sprintf("%v:%v", a.RemoteAddr, a.RemotePort)
		conn, err := net.Dial("udp", remote)
		share.Logger.Debug("new-client", zap.String("remote", remote), zap.String("local", conn.LocalAddr().String()))
		if err != nil {
			share.Logger.Error(remote, zap.String("err", fmt.Sprintf("%v", err)))
			time.Sleep(1 * time.Second)
			continue INITLOOP
		}

		share.Logger.Debug("whoami", zap.String("stat", "function started"), zap.String("IP", a.RemoteAddr), zap.Uint16("PORT", a.RemotePort))
		for z := 0; z < 3; z++ {
			if z == 0 {
				whoami, err = a.Whoami(&conn)
			} else if z == 3 {
				whoamiOld := whoami
				whoami, err = a.Whoami(&conn)
				if whoamiOld.IP != whoami.IP || whoamiOld.PORT != whoami.PORT {
					share.Logger.Error("whoami", zap.String("remote", remote), zap.String("error", "IP or Port is not stabile"))
				}
			} else {
				_, err = a.Whoami(&conn)
			}

			if err != nil {
				share.Logger.Error("whoami", zap.String("remote", remote), zap.Error(err))
				continue INITLOOP
			}
		}

		share.Logger.Info("whoami", zap.String("IP", whoami.IP), zap.String("PORT", whoami.PORT))
		share.Logger.Debug("connect", zap.String("stat", "function started"), zap.String("IP", a.RemoteAddr), zap.Uint16("PORT", a.RemotePort))
		err = a.Connect(&conn, clientId, whoami)
		if err != nil {
			share.Logger.Error("connect", zap.String("remote", remote), zap.Error(err))
			continue INITLOOP
		}
		share.Logger.Debug("connectDone", zap.String("IP", whoami.IP), zap.String("PORT", whoami.PORT))

		time.Sleep(time.Second * 5)

		if ScriptFile != "" {
			ex := share.Env{"MODE=client", "REMOTE_ADDR=" + a.RemoteAddr, "REMOTE_PORT=" + strconv.Itoa(int(a.RemotePort)), "MTU=" + strconv.Itoa(a.MTU), "REMOTE_LOCAL_IPV6=" + a.RemoteV6LocalAddr, "WHOAMI_IP=" + whoami.IP, "WHOAMI_PORT=" + whoami.PORT, "INTERFACE=" + share.InterfaceName(clientId), "FOU_PORT=" + strconv.Itoa(fouPortInt[clientId]), "REMOTE_FOU_PORT=" + strconv.Itoa(whoami.REMOTE_FOU_PORT)}
			stdout, stderr := ex.Exec(ScriptFile)
			share.Logger.Debug("script", zap.String("remote", remote), zap.String("stdout", stdout), zap.String("stderr", stderr))

		}

		err = a.Ping(clientId)
		if err != nil {
			share.Logger.Error("ping", zap.String("remote", remote), zap.Error(err))
		}
	}

}

func Shutdown(clients ClientConfig) int {
	for i := 0; i < len(clients.Servers); i++ {
		share.FouDel(fouPortInt[i+1])
		share.InterfaceDel(i + 1)
	}
	return 0
}

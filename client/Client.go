package client

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ahmetozer/dynamic-fou/share"
	"go.uber.org/zap"
)

var (
	ClientName string
)

func Start() {

	configFile := os.Getenv("CONFIG_FILE")

	if configFile == "" {
		configFile = "/etc/dynamic-fou.client.json"
	}
	share.Logger.Debug("Opening config file", zap.String("config-file", configFile))
	err := share.CheckFolder(filepath.Dir(configFile))
	if err != nil {
		share.Logger.Fatal("log path is not oppened", zap.String("path", filepath.Dir(configFile)), zap.String("err", err.Error()))
	}

	clients, err := Parse(configFile)
	if err != nil {
		share.Logger.Fatal("config file cannot parsed", zap.String("file", configFile), zap.String("err", err.Error()))
	}

	share.Logger.Debug("Config file parsed", zap.String("clientName", clients.ClientName), zap.String("server-count", fmt.Sprint((len(clients.Servers)))))
	ClientName = clients.ClientName
	var wg sync.WaitGroup
	for i := 0; i < len(clients.Servers); i++ {
		wg.Add(1)
		share.Logger.Debug(fmt.Sprintf("client %v", i+1), clients.Servers[i].toZap()...)
		go clients.Servers[i].clientInit()
	}

	share.Logger.Debug("client inits are done.")

	wg.Wait()
}

func (a SvConfig) clientInit() {
	k := true
	var IP, PORT string
INITLOOP:
	for k {
		remote := fmt.Sprintf("%v:%v", a.RemoteAddr, a.RemotePort)
		conn, err := net.Dial("udp", remote)
		share.Logger.Debug("new-client", zap.String("remote", remote), zap.String("local", conn.LocalAddr().String()))
		if err != nil {
			share.Logger.Error(remote, zap.String("err", fmt.Sprintf("%v", err)))
			continue INITLOOP
		}

		time.Sleep(1 * time.Second)
		share.Logger.Debug("whoami", zap.String("stat", "function started"), zap.String("IP", a.RemoteAddr), zap.Uint16("PORT", a.RemotePort))
		IP, PORT, err = a.Whoami(&conn)
		if err != nil {
			share.Logger.Error("whoami", zap.String("remote", remote), zap.Error(err))
			continue INITLOOP
		}
		share.Logger.Info("whoami", zap.String("IP", IP), zap.String("PORT", PORT))
		share.Logger.Debug("connect", zap.String("stat", "function started"), zap.String("IP", a.RemoteAddr), zap.Uint16("PORT", a.RemotePort))
		err = a.Connect(&conn)
		if err != nil {
			share.Logger.Error("connect", zap.String("remote", remote), zap.Error(err))
			continue INITLOOP
		}
		share.Logger.Info("done")

		time.Sleep(1 * time.Second)
	}

}

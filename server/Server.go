package server

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ahmetozer/dynamic-fou/share"
	"go.uber.org/zap"
)

var (
	configList []ClientConfig
	ClientList map[string]CurrentClient
)

type Config struct {
	PORT string
	IP   string
}

func Start() {
	PORT := os.Getenv("PORT")
	IP := os.Getenv("IP")
	if IP == "" {
		IP = "::"
	}

	if PORT == "" {
		PORT = "9000"
	}

	configFile := os.Getenv("CONFIG_FILE")

	if configFile == "" {
		configFile = "/etc/dynamic-fou.server.json"
	}

	share.Logger.Debug("Opening config file", zap.String("config-file", configFile))
	err := share.CheckFolder(filepath.Dir(configFile))
	if err != nil {
		share.Logger.Fatal("log path is not oppened", zap.String("path", filepath.Dir(configFile)), zap.String("err", err.Error()))
	}

	configList, err = Parse(configFile)
	if err != nil {
		share.Logger.Fatal("config file cannot parsed", zap.String("file", configFile), zap.String("err", err.Error()))
	}
	share.Logger.Debug("Config file parsed", zap.String("client-count", fmt.Sprint(len(configList))))

	for i := 0; i < len(configList); i++ {
		share.Logger.Debug(fmt.Sprintf("client %v", i+1), configList[i].toZap()...)
	}

	share.Logger.Debug("server info", zap.String("port", PORT), zap.String("ip", IP))

	i, err := strconv.Atoi(PORT)
	if err != nil {
		share.Logger.Fatal("PORT must be a number", zap.String("port", PORT), zap.String("err", err.Error()))
	}

	fouPort := os.Getenv("FOU_PORT")

	if fouPort == "" {
		fouPort = "5555"
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: i,
		IP:   net.ParseIP(IP),
	})
	if err != nil {
		share.Logger.Fatal(err.Error())
	}

	defer conn.Close()
	share.Logger.Info("server started", zap.String("server", conn.LocalAddr().String()))
	message := make([]byte, 2048)

StartLoop:
	for {

		rlen, remote, err := conn.ReadFromUDP(message[:])
		if err != nil {
			share.Logger.Error(err.Error())
			message = []byte{}
			continue StartLoop
		}
		go MessageTypeController(conn, message, rlen, remote)
	}
}

func MessageTypeController(conn *net.UDPConn, message []byte, rlen int, remote *net.UDPAddr) {
	data := strings.TrimSpace(string(message[:rlen]))

	cli, err := getClientByName(share.IniVal(data, "client"))
	if err != nil {
		share.Logger.Info("client info not parsed", zap.String("remote", remote.String()))
		return
	}
	otkStatus := share.OTKCheck(share.IniVal(data, "otk"), cli.ClientKey)
	mode := share.IniVal(data, "mode")
	share.Logger.Debug("newMesagge", zap.String("mode", mode), zap.String("otk", share.IniVal(data, "otk")), zap.String("remote", remote.String()), zap.Bool("isOTKValid", otkStatus))
	if !otkStatus {
		return
	}
	switch mode {
	case "whoami":
		Whoami(conn, remote)
	case "connect":
		Connect(conn, remote)
	default:
		share.Logger.Debug("unknow mode type", zap.String("remote", remote.String()), zap.String("mode", mode))
	}
}

func Whoami(conn *net.UDPConn, remote *net.UDPAddr) {
	_, err := conn.WriteTo([]byte(fmt.Sprintf("IP=%v\nPORT=%v\n", remote.IP, remote.Port)), remote)
	if err != nil {
		share.Logger.Error("whoami", zap.String("remote", remote.String()), zap.Error(err))
	}
}

func Connect(conn *net.UDPConn, remote *net.UDPAddr) {
	err := share.UDPPortCheck(remote.Port)
	status := "ok"
	if err != nil {
		status = "err"
		share.Logger.Debug("connect", zap.String("remote", remote.String()), zap.Error(err))
		_, err = conn.WriteTo([]byte(fmt.Sprintf("status=%v\nerr=%v\n", status, err)), remote)
		if err != nil {
			share.Logger.Error("connect", zap.String("remote", remote.String()), zap.Error(err))
		}
		status = "err"
		return
	}
	_, err = conn.WriteTo([]byte(fmt.Sprintf("status=%v\n", status)), remote)
	if err != nil {
		share.Logger.Error("connect", zap.String("remote", remote.String()), zap.Error(err))
	}

}

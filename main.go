package main

import (
	"os"
	"strconv"

	"github.com/ahmetozer/dynamic-fou/client"
	"github.com/ahmetozer/dynamic-fou/server"
	"github.com/ahmetozer/dynamic-fou/share"
	"go.uber.org/zap"
)

const (
	ProjectName string = "dynamic-fou"
)

var (
	printsdtout bool
	LOGLEVEL    string
	PORT        string
	IP          string
	MODE        string
	LOGFILE     string
)

func init() {
	PORT = os.Getenv("PORT")
	IP = os.Getenv("IP")
	MODE = os.Getenv("MODE")
	if MODE == "" {
		if len(os.Args) > 1 {
			MODE = os.Args[1]
		}
	}
	if MODE == "" {
		MODE = "server"
	}

	LOGFILE = os.Getenv("LOG_FILE")
	if LOGFILE == "" {
		LOGFILE = "-"
	}
	LOGLEVEL = os.Getenv("LOG_LEVEL")
	if LOGLEVEL == "" {
		LOGLEVEL = "2"
	}
}

func main() {
	u, err := strconv.ParseUint(LOGLEVEL, 10, 64)
	share.Err(err)
	// Start the logger
	share.InitLogger(LOGFILE, uint8(u))
	defer share.LogDefer()

	err = share.CheckKernelFouCapability()
	if err != nil && os.Getenv("KERNEL_FOU_TEST") != "no" {
		share.Logger.Fatal("sys-fou-test", zap.Error(err))
	}

	if MODE == "server" {
		share.Logger.Info("Starting server")
		server.Start()
	} else if MODE == "client" {
		share.Logger.Info("Starting client")
		client.Start()
	}

}

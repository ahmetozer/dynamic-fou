package main

import (
	"os"
	"strconv"

	"github.com/ahmetozer/dynamic-fou/client"
	"github.com/ahmetozer/dynamic-fou/server"
	"github.com/ahmetozer/dynamic-fou/share"
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
		MODE = "server"
	}
	LOGFILE = os.Getenv("LOG_FILE")
	if LOGFILE == "" {
		LOGFILE = "-"
	}
	LOGLEVEL = os.Getenv("LOG_LEVEL")
	if LOGLEVEL == "" {
		LOGLEVEL = "1"
	}
}

func main() {
	u, err := strconv.ParseUint(LOGLEVEL, 10, 64)
	share.Err(err)
	// Start the logger
	share.InitLogger(LOGFILE, uint8(u))
	defer share.LogDefer()

	if MODE == "server" {
		var a = server.Config{
			PORT: PORT,
			IP:   IP,
		}
		share.Logger.Info("Starting server")
		server.Start(&a)
	} else if MODE == "client" {
		share.Logger.Info("Starting client")
		client.Start()
	}

}

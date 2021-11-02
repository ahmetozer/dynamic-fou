package server

import (
	"net"
	"os"
	"strconv"

	"github.com/ahmetozer/dynamic-fou/share"
	"go.uber.org/zap"
)

func Pong(port int) {
	l, err := net.Listen("tcp", "[::]"+":"+strconv.Itoa(port))
	if err != nil {
		share.Logger.Panic("Pong", zap.Error(err))
		os.Exit(1)
	}
	share.Logger.Info("tcp conn test server started", zap.String("server", "[::]"+":"+strconv.Itoa(port)))
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			share.Logger.Error("Pong", zap.Error(err))
		}
		go PongHandleRequest(conn)
	}
}

// Handles incoming requests.
func PongHandleRequest(conn net.Conn) {
	share.Logger.Info("Pong", zap.String("newclient", conn.RemoteAddr().Network()))
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		share.Logger.Info("Pong", zap.String("client", conn.RemoteAddr().Network()), zap.Error(err))
	}

	conn.Close()
}

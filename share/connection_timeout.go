package share

import (
	"fmt"
	"net"
	"time"
)

func ConnectionTimeout(conn *net.Conn, s int, q *chan bool) {
	countDown := 0
	for {
		select {
		case <-*q:
			return
		default:
			if countDown < s {
				countDown += 1
				time.Sleep(time.Second)
			} else {
				fmt.Println("conection time out")
				(*conn).Close()
				return
			}
		}
	}
}

package share

import (
	"crypto/md5"
	"fmt"
	"time"
)

func NewOTK(password string) string {
	now := time.Now()
	secs := now.Unix()

	data := []byte(fmt.Sprintf("%v%v", password, secs/30))

	return fmt.Sprintf("%x", md5.Sum(data))
}

func OTKCheck(otk, password string) bool {
	return otk == NewOTK(password)
}

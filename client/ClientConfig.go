package client

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ClientConfig struct {
	ClientName string
	Servers    []SvConfig
}
type SvConfig struct {
	RemotePort      uint16
	ControlInterval uint16
	ControlRetry    uint16
	RemoteAddr      string
	ServerKey       string
	ClientKey       string
}

func Parse(configPath string) (ClientConfig, error) {
	config, err := os.ReadFile(configPath)
	if err != nil {
		return ClientConfig{}, err
	}

	dataJson := string(config)
	var arr ClientConfig
	err = json.Unmarshal([]byte(dataJson), &arr)
	if err != nil {
		return ClientConfig{}, err
	}
	return arr, nil
}

// return server config as string
func (k SvConfig) toString() string {
	v := reflect.ValueOf(k)
	typeOfS := v.Type()

	var t string
	for i := 0; i < v.NumField(); i++ {
		t += fmt.Sprintf("\"%s\":\"%v\",", typeOfS.Field(i).Name, v.Field(i).Interface())
	}
	return t
}

// return server config as zapcoreField
func (k SvConfig) toZap() []zapcore.Field {
	v := reflect.ValueOf(k)
	typeOfS := v.Type()

	var r []zapcore.Field
	for i := 0; i < v.NumField(); i++ {
		r = append(r, zap.String(typeOfS.Field(i).Name, fmt.Sprintf("%v", v.Field(i).Interface())))
	}
	return r
}

package server

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
	ServerKey  string
	ClientKey  string
	Route      []string
	Route6     []string
	Addr       []string
	Addr6      []string
}

func Parse(configPath string) ([]ClientConfig, error) {
	config, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	dataJson := string(config)
	var arr []ClientConfig
	err = json.Unmarshal([]byte(dataJson), &arr)
	if err != nil {
		return nil, err
	}
	return arr, nil
}

// return server config as string
func (k ClientConfig) toString() string {
	v := reflect.ValueOf(k)
	typeOfS := v.Type()

	var t string
	for i := 0; i < v.NumField(); i++ {
		t += fmt.Sprintf("\"%s\":\"%v\",", typeOfS.Field(i).Name, v.Field(i).Interface())
	}
	return t
}

// return server config as zapcoreField
func (k ClientConfig) toZap() []zapcore.Field {
	v := reflect.ValueOf(k)
	typeOfS := v.Type()

	var r []zapcore.Field
	for i := 0; i < v.NumField(); i++ {
		r = append(r, zap.String(typeOfS.Field(i).Name, fmt.Sprintf("%v", v.Field(i).Interface())))
	}
	return r
}

func getClientByName(name string) (ClientConfig, error) {
	for i := 0; i < len(configList); i++ {
		if configList[i].ClientName == name {
			return configList[i], nil
		}
	}
	return ClientConfig{}, fmt.Errorf("client not found")
}

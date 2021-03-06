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
	MTU        int
	ClientName string
	ClientKey  string
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
	l := v.NumField() - 1
	for i := 0; i < l+1; i++ {
		if i < l {
			t += fmt.Sprintf("\"%s\":\"%v\",", typeOfS.Field(i).Name, v.Field(i).Interface())
		} else {
			t += fmt.Sprintf("\"%s\":\"%v\"", typeOfS.Field(i).Name, v.Field(i).Interface())
		}
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

type CurrentClient struct {
	IP   string
	PORT int
}

// return server config as string
func (k CurrentClient) toString() string {
	v := reflect.ValueOf(k)
	typeOfS := v.Type()

	var t string

	l := v.NumField() - 1
	for i := 0; i < l+1; i++ {
		if i < l {
			t += fmt.Sprintf("\"%s\":\"%v\",", typeOfS.Field(i).Name, v.Field(i).Interface())
		} else {
			t += fmt.Sprintf("\"%s\":\"%v\"", typeOfS.Field(i).Name, v.Field(i).Interface())
		}

	}
	return t
}

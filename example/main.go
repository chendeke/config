package main

import (
	"fmt"
	"github.com/chendeke/config"
	"github.com/chendeke/config/config/reader"
)

type WatchConfig struct {
	Level        string                 `json:"level" default:"debug"`
	Encode       string                 `json:"encode"`
	LevelPort    int                    `json:"level_port" default:"9090"`
	LevelPattern string                 `json:"level_pattern" default:"/handle/okokok"`
	InitFields   map[string]interface{} `json:"init_fields"`
}

func main() {
	wc := new(WatchConfig)
	err := config.Get("logs").Scan(wc)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(*wc)
	}

	exit := make(chan struct{})

	changeLogLevel := func(value reader.Value) {
		fmt.Println("main:", string(value.Bytes()))
	}

	if watchObj, err := config.WatchFunc(changeLogLevel, "logs.level"); err == nil {
		defer func() { _ = watchObj.Close() }()
	} else {
		panic(err.Error())
	}

	// time.Sleep(time.Second)
	<-exit
}

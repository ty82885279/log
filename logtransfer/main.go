package main

import (
	"fmt"
	"logtransfer/config"
	"logtransfer/es"
	"logtransfer/kafka"

	"gopkg.in/ini.v1"
)

var (
	cfg = new(config.Config)
)

func main() {
	fmt.Println("hello")
	err := ini.MapTo(cfg, "./config/config.ini")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", cfg)
	//初始化ES
	es.Init(cfg.Es.Address)
	//从kafka中消费
	kafka.ReadMsgFromKafka(cfg.Kafka.Address, cfg.Kafka.Topic)
	select {}
}

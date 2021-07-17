package main

import (
	"fmt"
	"logagent/config"
	"logagent/etcd"
	"logagent/kafka"
	taillog "logagent/tailf"
	"logagent/utils"
	"sync"

	"gopkg.in/ini.v1"

	"time"
)

var (
	cfg = new(config.App)
)

// logAgent 程序入口
func main() {
	// 0.加载配置文件
	err := ini.MapTo(cfg, "./config/config.ini")
	if err != nil {

		fmt.Printf("mapping cfg failed,err:%#v\n", err)
		return
	}
	// 1.初始化kafka连接
	address := []string{cfg.Kafka.Address}
	err = kafka.Init(address, cfg.Kafka.MaxSize)
	if err != nil {
		fmt.Println("connect kafka failed")
		return
	}
	fmt.Println("Init kafka success")

	// 2. 从ETCD里拉去最新的配置项
	err = etcd.Init(cfg.Etcd.Address, time.Duration(cfg.Etcd.Timeout)*time.Second)
	if err != nil {
		fmt.Printf("Connect etcd failed,err:%#v\n", err)
		return
	}
	etctKey := fmt.Sprintf(cfg.Etcd.Key, utils.GetPulicIP())
	fmt.Println("Init etcd success")
	logEntries, err := etcd.GetConfig(etctKey)
	if err != nil {
		fmt.Printf("get cfg from ectd failed,err:%#v\n", err)
		return
	}

	// 3.打开日志文件准备收集日志
	taillog.Init(logEntries)
	configChan := taillog.NewConfigChan()

	// 派一个哨兵监视key的变化
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConf(etctKey, configChan)
	wg.Wait()

	// 3.1 循环每一个日志收集项 创建tailObj
	// 3.2
	//run()
	/*

	 */

}

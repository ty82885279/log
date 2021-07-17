package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"logagent/config"
	"time"

	"go.etcd.io/etcd/clientv3"
)

var (
	cli *clientv3.Client
)

// 初始化 etcd
func Init(address string, timeout time.Duration) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{address},
		DialTimeout: timeout,
	})
	if err != nil {
		//fmt.Printf("connect etcd failed,%#v\n", err)
		return
	}
	return
}

//获取最新的配置
func GetConfig(key string) (LogEntryConf []*config.LogEntry, err error) {
	// get
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)

		return
	}
	for _, ev := range resp.Kvs {

		err = json.Unmarshal(ev.Value, &LogEntryConf)
		if err != nil {
			fmt.Printf("json unmarshal failed,err:%#v\n", err)
			return
		}
	}

	return LogEntryConf, err
}

func WatchConf(key string, newConfChan chan<- []*config.LogEntry) {
	rch := cli.Watch(context.Background(), key) // <-chan WatchResponse
	for wresp := range rch {
		for _, ev := range wresp.Events {
			//fmt.Println("新配置来了")
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value) //在这里监听配置变化，发送到通道中
			var configs []*config.LogEntry
			if ev.Type != clientv3.EventTypeDelete {

				err := json.Unmarshal(ev.Kv.Value, &configs)

				if err != nil {
					fmt.Printf("unmarshal failed,err:%#v\n", err)
				}
			}

			newConfChan <- configs
		}
	}
}

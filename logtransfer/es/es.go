package es

import (
	"context"
	"fmt"
	"logtransfer/config"
	"time"

	"github.com/olivere/elastic/v7"
)

var (
	ESClient *elastic.Client
	DataCh   = make(chan *config.LogData, 100000)
)

func Init(addr string) {
	var err error
	ESClient, err = elastic.NewClient(elastic.SetURL("http://" + addr))
	if err != nil {
		fmt.Printf("connect es failed,err:", err)
		return
	}

	fmt.Println("connect es success")
	//ch = make(chan *config.LogData, 100000)
	for i := 0; i < 16; i++ {
		go sendToEsFromCh()
	}
	return
}

func SendData(ld *config.LogData) {
	DataCh <- ld
}

// 从通道中取出数据发送到ES中
func sendToEsFromCh() {
	for {
		select {
		case ld := <-DataCh:
			put1, err := ESClient.Index().
				Index(ld.Topic).
				BodyJson(ld).
				Do(context.Background())
			if err != nil {
				fmt.Printf("put es failed,err:", err)
			}
			fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
		default:
			time.Sleep(500 * time.Millisecond)
		}
	}
}

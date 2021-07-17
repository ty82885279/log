package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

//专门往kafka写日志

type logData struct {
	topic string
	data  string
}

var (
	client      sarama.SyncProducer
	logDataChan chan *logData
)

func Init(address []string, maxSize int) (err error) {
	//fmt.Println("init kafka")

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	// 连接kafka
	client, err = sarama.NewSyncProducer(address, config)

	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	logDataChan = make(chan *logData, maxSize)
	go sendToKafka()
	return
}
func SendChan(topic, data string) {

	msg := &logData{
		topic: topic,
		data:  data,
	}
	logDataChan <- msg
}

func sendToKafka() {
	for {
		select {
		case ld := <-logDataChan:
			//构建消息
			msg := &sarama.ProducerMessage{}
			msg.Topic = ld.topic
			msg.Value = sarama.StringEncoder(ld.data)

			//发送消息
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("send msg failed, err:", err)
				return
			}
			fmt.Printf("pid:%v offset:%v\n", pid, offset)
			fmt.Printf("topic:%v value:%v\n", ld.topic, ld.data)
		default:
			time.Sleep(time.Millisecond * 50)

		}

	}

}

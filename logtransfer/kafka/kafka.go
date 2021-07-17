package kafka

import (
	"fmt"
	"logtransfer/config"
	"logtransfer/es"

	"github.com/Shopify/sarama"
)

var (
	consumer sarama.Consumer
)

func ReadMsgFromKafka(addr string, topic string) {
	var err error
	consumer, err = sarama.NewConsumer([]string{addr}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println(partitionList)
	fmt.Println("长度--", len(partitionList))

	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		//defer pc.AsyncClose() //不能关闭
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:-%v\n",
					msg.Partition, msg.Offset, msg.Key, string(msg.Value))
				//向ES发送日志
				var ld = config.LogData{Topic: topic, Data: string(msg.Value)}
				es.SendData(&ld) //数据放到缓冲通道中立即返回
			}
		}(pc)

	}
}

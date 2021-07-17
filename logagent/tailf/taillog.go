package taillog

import (
	"context"
	"fmt"
	"logagent/kafka"

	"github.com/hpcloud/tail"
)

//一个日志收集的任务
type TailTask struct {
	path       string
	topic      string
	instance   *tail.Tail
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewTailTask(path, topic string) (tailObj *TailTask) {
	ctx, cancel := context.WithCancel(context.Background())
	tailObj = &TailTask{
		path:       path,
		topic:      topic,
		ctx:        ctx,
		cancelFunc: cancel,
	}
	tailObj.init()

	return tailObj
}

func (t TailTask) init() {
	conf := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	}
	var err error
	t.instance, err = tail.TailFile(t.path, conf)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
		return
	}
	go t.run() //直接去采集日志发送到kafka
}

func (t *TailTask) run() {
	for {
		select {
		case <-t.ctx.Done():
			fmt.Printf("任务已经取消：%s_%s\n", t.path, t.topic)
			return
		case line := <-t.instance.Lines:
			//kafka.SendToKafka(t.topic, line.Text)
			kafka.SendChan(t.topic, line.Text) //放到缓存通道中，立即返回
		}
	}
}

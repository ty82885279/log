package taillog

import (
	"fmt"
	"logagent/config"
	"time"
)

var (
	tskMgr *tailLogMgr
)

type tailLogMgr struct {
	logEntry    []*config.LogEntry
	taskMap     map[string]*TailTask
	newConfChan chan []*config.LogEntry
}

func Init(LogEntryConf []*config.LogEntry) {
	tskMgr = &tailLogMgr{
		logEntry:    LogEntryConf,
		taskMap:     make(map[string]*TailTask, 16),
		newConfChan: make(chan []*config.LogEntry),
	}
	for _, logEntry := range LogEntryConf {
		task := NewTailTask(logEntry.Path, logEntry.Topic)
		key := fmt.Sprintf("%s%s", logEntry.Path, logEntry.Topic)
		tskMgr.taskMap[key] = task
	}

	go tskMgr.run()
}

func (t tailLogMgr) run() {
	var tempConf []*config.LogEntry
	for {
		select {
		case newConf := <-t.newConfChan:
			tempConf = newConf
			fmt.Println("Config has changed")
			for _, value := range newConf {
				key := fmt.Sprintf("%s%s", value.Path, value.Topic)
				_, ok := t.taskMap[key]
				if !ok {
					tailTask := NewTailTask(value.Path, value.Topic)
					t.taskMap[key] = tailTask
				} else {
					fmt.Println("This task already exists")
					continue
				}
			}
			//从原来的配置中和新的配置挨个对比，如果都相同，不做操作。不同，删除原来的配置项（标识位判断）,增加2个字段，ctx.cancelFunc,
			for _, c1 := range t.logEntry {
				isDelete := true //删除配置的标志位
				for _, c2 := range newConf {
					if c1.Path == c2.Path && c1.Topic == c2.Topic {
						isDelete = false //不需要删除
						continue
					}
				}
				if isDelete {
					//把c1对应的后台任务取消
					key := fmt.Sprintf("%s%s", c1.Path, c1.Topic)
					task := t.taskMap[key]
					task.cancelFunc()
					delete(t.taskMap, key)
				}
			}
			t.logEntry = tempConf
		default:
			time.Sleep(time.Second)
		}
	}
}

func NewConfigChan() chan<- []*config.LogEntry {

	return tskMgr.newConfChan
}

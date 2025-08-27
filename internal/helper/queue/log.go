package queue

import (
	"context"
	"time"

	table "cdp-admin-service/internal/helper/dal"

	"github.com/zeromicro/go-zero/core/logx"
)

type Queue struct {
	queue chan table.TCdpSysLog
}

func InitLogQueue() *Queue {
	logQueue := &Queue{
		queue: make(chan table.TCdpSysLog, 100),
	}
	go logQueue.StartConsume()
	return logQueue
}

func (l *Queue) Producer(data table.TCdpSysLog) {
	select {
	case l.queue <- data:
	default:
		// 队列满，丢弃
		logx.Errorf("log queue is full, drop data[%+v]", data)
		return
	}
}
func (l *Queue) StartConsume() {

	for {
		data := <-l.queue
		table.T_TCdpSysLogService.Insert(context.Background(), "cdp_sys_log", &data)
		time.Sleep(5 * time.Millisecond)
	}
}

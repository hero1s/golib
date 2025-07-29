package toolkit

import (
	"github.com/hero1s/golib/log"
	"time"
)

// 简单计时器
type ClockLog struct {
	name      string
	startTime time.Time
	timeout   time.Duration
}

// NewSimpleTimer 创建简单计时器
func NewClockLog(name string, timeout time.Duration) *ClockLog {
	st := &ClockLog{
		name:      name,
		startTime: time.Now(),
		timeout:   timeout,
	}
	log.Infof("[%s] 计时开始 at %v\n", st.name, st.startTime.Format("15:04:05.000"))
	return st
}

// End 结束计时并输出日志
func (st *ClockLog) End() {
	endTime := time.Now()
	duration := endTime.Sub(st.startTime)
	if duration > st.timeout {
		log.Errorf("[%s] 计时结束 at %v, 耗时: %v, 超时: %v\n", st.name, endTime.Format("15:04:05.000"), duration, st.timeout)
	} else {
		log.Infof("[%s] 计时结束 at %v, 耗时: %v\n", st.name, endTime.Format("15:04:05.000"), duration)
	}
}

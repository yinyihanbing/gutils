package goutils

import (
	"github.com/yinyihanbing/gutils/timer"
	"time"
)

// 定时器结构体
type TimerHelper struct {
	dispatcher *timer.Dispatcher
	closeSig   chan bool
}

// 实例定时器帮助类
func NewTimerHelper() *TimerHelper {
	timerCronExpr := new(TimerHelper)
	timerCronExpr.dispatcher = timer.NewDispatcher(10000)
	timerCronExpr.closeSig = make(chan bool, 1)

	go timerCronExpr.run()

	return timerCronExpr
}

// 添加定时任务
// 示例:
// 每隔5秒执行一次：*/5 * * * * *
// 每隔1分钟执行一次：0 */1 * * * *
// 每天23点执行一次：0 0 23 * * *
// 每天凌晨1点执行一次：0 0 1 * * *
// 每月1号凌晨1点执行一次：0 0 1 1 * *
// 在26分、29分、33分执行一次：0 26,29,33 * * * *
// 每天的0点、13点、18点、21点都执行一次：0 0 0,13,18,21 * * *
// 每周5凌晨0点执行 0 0 0 * * 5
func (this *TimerHelper) CronFuncExt(expr string, cb func()) *timer.Cron {
	cronExpr, err := timer.NewCronExpr(expr)
	if err != nil {
		panic("invalid CronExpr")
	}
	return this.dispatcher.CronFunc(cronExpr, cb)
}

// 等待指定时间后执行1次
func (this *TimerHelper) AfterFunc(d time.Duration, cb func()) *timer.Timer {
	return this.dispatcher.AfterFunc(d, cb)
}

// 执行定时函数
func (this *TimerHelper) run() {
	for {
		select {
		case <-this.closeSig:
			return
		case t := <-this.dispatcher.ChanTimer:
			t.Cb()
		}
	}
}

// 停止定时器
func (this *TimerHelper) Stop() {
	this.closeSig <- true
}

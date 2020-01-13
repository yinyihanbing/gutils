package goutils

import "sync"

// 时间轮
type TimingWheel struct {
	mu            sync.RWMutex
	Wheel         [][]interface{}   // 轮子
	Interval      int               // 定时频率间隔(秒)
	WheelSize     int               // 轮子大小
	FirstOffset   int               // 第一次运行隔几个偏移（新加任务隔几个偏移后执行）
	CurrentOffset int               // 当前执行偏移数
	WheelHandler  func(interface{}) // 轮子触发回调
}

// 新建时间轮
func NewTimeWheel(wheelSize, interval, firstOffset int, wheelHandler func(obj interface{})) *TimingWheel {
	sw := TimingWheel{}
	sw.WheelSize = wheelSize
	sw.Wheel = make([][]interface{}, sw.WheelSize)
	sw.Interval = interval
	sw.FirstOffset = firstOffset
	sw.WheelHandler = wheelHandler

	return &sw
}

// 添加轮子中的触发对象
func (this *TimingWheel) AddWheelObject(obj interface{}, offset int) {
	this.mu.RLock()
	defer this.mu.RUnlock()

	idx := (this.FirstOffset + this.CurrentOffset + offset) % this.WheelSize
	if this.Wheel[idx] == nil {
		this.Wheel[idx] = make([]interface{}, 0, 1)
	}
	this.Wheel[idx] = append(this.Wheel[idx], obj)
}

// 移除轮子中的触发对象
func (this *TimingWheel) RemoveWheelObject(obj interface{}, compare func(src, tar interface{}) bool) {
	this.mu.RLock()
	defer this.mu.RUnlock()

	exists := false
	for i, v1 := range this.Wheel {
		for j, v2 := range v1 {
			if compare(v2, obj) {
				this.Wheel[i] = append(this.Wheel[i][0:j], this.Wheel[i][j+1:]...)
				exists = true
				break
			}
		}
		if exists {
			break
		}
	}
}

// 轮子转动
func (this *TimingWheel) WheelTurn() {
	this.CurrentOffset += 1
	this.CurrentOffset = this.CurrentOffset % this.WheelSize
	if len(this.Wheel[this.CurrentOffset]) == 0 {
		return
	}
	for _, v := range this.Wheel[this.CurrentOffset] {
		this.WheelHandler(v)
	}
}

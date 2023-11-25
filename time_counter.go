package timecounter

import (
	"sync"
	"time"
)

type Slot struct {
	value float64
	gen   int64 // 周目数
	lock  sync.Mutex
}

type Counter struct {
	slots    []Slot
	span     int64
	slotNum  int64
	slotSpan int64
}

const DEFAULT_SLOT_NUM = 10

func NewCounter(span time.Duration, slotNum int64) *Counter {
	if slotNum <= 0 {
		slotNum = DEFAULT_SLOT_NUM
	}
	return &Counter{
		slots:    make([]Slot, slotNum),
		span:     span.Milliseconds(),
		slotNum:  slotNum,
		slotSpan: span.Milliseconds() / slotNum,
	}
}

func (cl *Counter) PutNow() {
	cl.PutN(time.Now(), 1)
}
func (cl *Counter) PutNowN(n float64) {
	cl.PutN(time.Now(), n)
}

func (cl *Counter) Put(now time.Time) {
	cl.PutN(now, 1)
}
func (cl *Counter) PutN(now time.Time, n float64) {
	totalSpan := now.UnixMilli()
	gen := totalSpan / cl.span                     // 存储为gen周目
	slotIndex := totalSpan % cl.span / cl.slotSpan // 存储到第slotIndex个slot上
	cl.slots[slotIndex].lock.Lock()
	if cl.slots[slotIndex].gen != gen {
		cl.slots[slotIndex].gen = gen
		cl.slots[slotIndex].value = 0
	}
	cl.slots[slotIndex].value += n
	cl.slots[slotIndex].lock.Unlock()
}

func (cl *Counter) SumNow(duration time.Duration) float64 {
	return cl.Sum(time.Now(), duration)
}

func (cl *Counter) SumAll() float64 {
	return cl.Sum(time.Now(), time.Duration(cl.span)*time.Millisecond)
}

func (cl *Counter) Sum(now time.Time, duration time.Duration) float64 {
	var sum float64
	var totalSpan = now.UnixMilli()
	gen := totalSpan / cl.span
	nowIndex := totalSpan%cl.span/cl.slotSpan - 1     // 当前所在slot位置的上一个slot
	slotsNum := duration.Milliseconds() / cl.slotSpan // 需要计算的slot数量
	for i := int64(0); i < cl.slotNum && i < slotsNum; i++ {
		index := (nowIndex - i + cl.slotNum) % cl.slotNum // 第i个slot的位置
		// 只有最近一个连续周目的数据才有效
		if diff := gen - cl.slots[index].gen; diff <= 1 && diff >= 0 {
			sum += cl.slots[index].value
		}
	}
	return sum
}

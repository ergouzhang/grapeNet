// 计算一个周期内数据是多少，超过这个周期会直接清理数据
// 主要用于统计、限制以及周期性的内存计算
// version 1.0 beta
// by koangel
// email: jackliu100@gmail.com
// 2021/01/09
package continer

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type timeCount struct {
	count    int32
	nextTime int64
	hotTime  int64
}

func (c *timeCount) tickAdd(limit int, t time.Duration) bool {
	if c.IsExpired() {
		atomic.StoreInt32(&c.count, 0)
		c.next(t)
	}

	if c.IsExpired() == false && c.count >= int32(limit) {
		return false
	}

	c.hotTime = time.Now().Add(40 * time.Minute).Unix()
	atomic.AddInt32(&c.count, 1)
	return true
}

func (c *timeCount) next(t time.Duration) {
	c.nextTime = time.Now().Add(t).Unix()
}

func (c *timeCount) IsExpired() bool {
	if time.Now().Unix() >= c.nextTime {
		return true
	}

	return false
}

func (c *timeCount) IsNotHot() bool {
	if time.Now().Unix() >= c.hotTime {
		return true
	}

	return false
}

type TimeGroup struct {
	mux      sync.RWMutex
	mapData  map[interface{}]*timeCount
	limit    int
	lootTime time.Duration
	once     *timeCount

	ticks  *time.Ticker
	ctx    context.Context
	cancel context.CancelFunc

	stopOnce sync.Once
}

func NewTimeGroup(loopTime time.Duration, clearTick time.Duration, limit int) *TimeGroup {
	ret := &TimeGroup{
		mapData:  map[interface{}]*timeCount{},
		limit:    limit,
		lootTime: loopTime,
		once:     nil,
		ticks:    time.NewTicker(clearTick),
	}

	ret.ctx, ret.cancel = context.WithCancel(context.Background())
	go ret.procClearup()

	return ret
}

func (c *TimeGroup) procClearup() {
	for {
		select {
		case <-c.ticks.C:
			{
				c.ClearMemory()
			}
		case <-c.ctx.Done():
			{
				return
			}
		}
	}
}

func (t *TimeGroup) Stop() {
	t.stopOnce.Do(func() {
		t.ClearMemory()
		t.cancel()
	})
}

func (t *TimeGroup) AddLimitMap(key interface{}) bool {
	if t.once == nil {
		c := &timeCount{
			count:    0,
			nextTime: 0,
		}
		c.next(t.lootTime)
		t.once = c
	}

	if t.once.IsExpired() {
		t.mux.Lock()
		t.mapData = map[interface{}]*timeCount{}
		t.mux.Unlock()
		t.once.next(t.lootTime)
	}

	if t.once.IsExpired() == false && len(t.mapData) >= t.limit {
		return false
	}

	t.mux.RLock()
	_, ok := t.mapData[key]
	t.mux.RUnlock()
	if !ok {
		t.mux.Lock()
		t.mapData[key] = &timeCount{
			count:    0,
			nextTime: 0,
		}
		t.mux.Unlock()
	}

	return true
}

func (t *TimeGroup) ClearMemory() {
	delKey := []interface{}{}
	t.mux.RLock()
	for k, v := range t.mapData {
		if v.IsNotHot() {
			delKey = append(delKey, k)
		}
	}
	t.mux.RUnlock()

	if len(delKey) == 0 {
		return
	}

	t.mux.Lock()
	for _, v := range delKey {
		delete(t.mapData, v)
	}
	t.mux.Unlock()
}

func (t *TimeGroup) AddCount(key interface{}) bool {

	t.mux.RLock()
	mv, ok := t.mapData[key]
	t.mux.RUnlock()
	if !ok {
		c := &timeCount{
			count:    0,
			nextTime: 0,
		}

		t.mux.Lock()
		t.mapData[key] = c
		t.mux.Unlock()

		c.next(t.lootTime)
		return c.tickAdd(t.limit, t.lootTime)
	}

	return mv.tickAdd(t.limit, t.lootTime)
}

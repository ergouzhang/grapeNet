// 定时执行函数的的数据容器
// version 1.0 beta
// by koangel
// email: jackliu100@gmail.com
// 2021/01/09
package continer

import (
	"context"
	"hash/crc32"
	"log"
	"sync"
	"time"
)

// 返回true 则自动清理数据 如果返回false则不清理
type TickWriter func(c *TickLimit, key string, data []interface{}) bool
type ExecProc func(c *TickLimit)
type limitType = map[string]*LimitAarray

type TickLimit struct {
	lc       int
	name     string
	limitArr []*LimitAarray

	ticker *time.Ticker

	startProc ExecProc
	caller    TickWriter
	endProc   ExecProc

	once     sync.Once
	cxtClose context.CancelFunc
	ctx      context.Context

	rl sync.RWMutex

	Cli        interface{}
	DebugPrint bool

	arrayCount int
	index      int
}

func GetHostIndex(key string, maxCount int) int {
	hashCode := crc32.ChecksumIEEE([]byte(key))
	return int(hashCode) % maxCount
}

func NewTickerLimit(name string, index int, arrCount, limit int, ticker time.Duration, fn TickWriter, start, end ExecProc) *TickLimit {
	r := &TickLimit{
		arrayCount: arrCount,
		index:      index,
		name:       name,
		lc:         limit,
		limitArr:   []*LimitAarray{},
		ticker:     time.NewTicker(ticker),
		caller:     fn,
		startProc:  start,
		endProc:    end,
		once:       sync.Once{},
		cxtClose:   nil,
		ctx:        nil,
		DebugPrint: false,
	}

	r.StartProc()
	return r
}

func (c *TickLimit) Name() string {
	return c.name
}

func (c *TickLimit) Index() int {
	return c.index
}

func (c *TickLimit) find(name string) *LimitAarray {
	c.rl.RLock()
	for i := 0; i < len(c.limitArr); i++ {
		if c.limitArr[i].name == name {
			r := c.limitArr[i]
			c.rl.RUnlock()
			return r
		}
	}
	c.rl.RUnlock()

	nl := NewLA(c.lc, name)
	c.rl.Lock()
	c.limitArr = append(c.limitArr, nl)
	c.rl.Unlock()

	return nl
}

func (c *TickLimit) Add(key string, data interface{}) {
	l := c.find(key)
	if l != nil {
		l.Add(data)
	}
}

func (c *TickLimit) StartProc() {
	c.ctx, c.cxtClose = context.WithCancel(context.Background())
	c.once.Do(func() {
		go c.onTicker()
	})
}

func (c *TickLimit) Stop() {
	c.cxtClose() // 断开
}

func (c *TickLimit) onTicker() {
	for {
		select {
		case <-c.ctx.Done():

			{
				return
			}
		case <-c.ticker.C:
			{
				if len(c.limitArr) <= 0 {
					continue
				}

				c.rl.Lock()
				arrData := c.limitArr
				c.limitArr = []*LimitAarray{}
				c.rl.Unlock()

				startTime := time.Now()
				procCount := 0
				c.startProc(c)
				for _, limit := range arrData {
					procCount += limit.Len()
					if c.caller(c, limit.Name(), limit.Array()) {
						limit.Clear()
					}
				}
				c.endProc(c)
				if c.DebugPrint {
					// debug for exec time...
					log.Println("tick run time:", time.Now().Sub(startTime))
				}
			}
		}
	}
}

func (c *TickLimit) Data(key string) []interface{} {
	value := []interface{}{}

	l := c.find(key)
	if l != nil {
		l.Foreach(func(n interface{}) {
			value = append(value, n)
		})
	}

	return value
}

func (c *TickLimit) Len(key string) int {

	l := c.find(key)
	if l != nil {
		return l.Len()
	}

	return 0
}

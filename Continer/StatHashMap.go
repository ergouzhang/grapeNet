package continer

// 内存数值统计用HASHMAP
// version 1.0 beta
// by koangel
// email: jackliu100@gmail.com
// 2020/10/10
import "sync"

type SStatMap struct {
	data   map[string]int64
	locker sync.RWMutex
}

func NewStatMap() *SStatMap {
	return &SStatMap{
		data: map[string]int64{},
	}
}

func (c *SStatMap) Len() int {
	c.locker.RLock()
	defer c.locker.RUnlock()

	return len(c.data)
}

func (c *SStatMap) Value(key string) int64 {
	c.locker.RLock()
	defer c.locker.RUnlock()

	value, ok := c.data[key]
	if !ok {
		return 0
	}

	return value
}

func (c *SStatMap) Incr(key string, value int64) {
	c.locker.Lock()
	defer c.locker.Unlock()

	val, ok := c.data[key]
	if ok {
		c.data[key] = val + value
	} else {
		c.data[key] = value
	}
}

func (c *SStatMap) Sub(key string, value int64) {
	c.locker.Lock()
	defer c.locker.Unlock()

	val, ok := c.data[key]
	if ok {
		c.data[key] = val - value
	} else {
		c.data[key] = 0 - value
	}
}

func (c *SStatMap) Range(fn func(key string, value int64) bool) {
	c.locker.RLock()
	defer c.locker.RUnlock()

	for key, value := range c.data {
		if fn(key, value) == false {
			break
		}
	}
}

func (c *SStatMap) Clear() {
	c.locker.Lock()
	defer c.locker.Unlock()

	c.data = map[string]int64{}
}

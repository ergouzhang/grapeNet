package continer

import "sync"

// 2021/01/09 新增高性能的容器，LimitArray，采用数组切片完成计算，效率极高
type LimitAarray struct {
	limit  int
	con    []interface{}
	locker sync.RWMutex
	name   string
}

func NewLA(lc int, key string) *LimitAarray {
	return &LimitAarray{
		limit: lc,
		name:  key,
		con:   []interface{}{},
	}
}

func (l *LimitAarray) Name() string {
	return l.name
}

func (l *LimitAarray) Add(val interface{}) {
	l.locker.Lock()
	defer l.locker.Unlock()

	l.con = append(l.con, val)

	// 循环删除
	if len(l.con) > l.limit {
		starPos := len(l.con) - l.limit
		l.con = l.con[starPos:]
	}
}

func (l *LimitAarray) Array() []interface{} {
	return l.con
}

func (l *LimitAarray) Clear() {
	l.locker.Lock()
	defer l.locker.Unlock()

	l.con = []interface{}{}
}

func (l *LimitAarray) RevForeach(fn func(n interface{})) {
	l.locker.RLock()
	defer l.locker.RUnlock()

	for i := len(l.con) - 1; i >= 0; i-- {
		fn(l.con[i])
	}
}

func (l *LimitAarray) Foreach(fn func(n interface{})) {
	l.locker.RLock()
	defer l.locker.RUnlock()

	for i := 0; i < len(l.con); i++ {
		fn(l.con[i])
	}
}

func (l *LimitAarray) Len() int {
	return len(l.con)
}

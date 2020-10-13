// 一款带锁的列表
// version 1.0 beta
// by koangel
// email: jackliu100@gmail.com
// 2017/12/10
package continer

import (
	"container/list"
	"sync"
)

type SList struct {
	slist  *list.List
	locker sync.RWMutex
}

func New() *SList {
	return &SList{
		slist: list.New(),
	}
}

func (sc *SList) Push(item interface{}) {
	sc.locker.Lock()
	defer sc.locker.Unlock()

	sc.slist.PushBack(item)
}

func (sc *SList) First() interface{} {
	sc.locker.RLock()
	defer sc.locker.RUnlock()

	return sc.slist.Front()
}

func (sc *SList) Back() interface{} {
	sc.locker.RLock()
	defer sc.locker.RUnlock()

	return sc.slist.Back()
}

func (sc *SList) Clear() {
	sc.locker.Lock()
	defer sc.locker.Unlock()

	sc.slist = list.New() //
}

func (sc *SList) Range(fn func(i interface{}) bool) {
	sc.locker.RLock()
	defer sc.locker.RUnlock()

	for e := sc.slist.Front(); e != nil; e = e.Next() {
		if fn(e.Value) == false {
			break
		}
	}
}

func (sc *SList) ReverseRange(fn func(i interface{}) bool) {
	sc.locker.RLock()
	defer sc.locker.RUnlock()

	for e := sc.slist.Back(); e != nil; e = e.Prev() {
		if fn(e.Value) == false {
			break
		}
	}
}

func (sc *SList) Search(fn func(i interface{}) bool) (interface{}, bool) {
	sc.locker.RLock()
	defer sc.locker.RUnlock()

	for e := sc.slist.Front(); e != nil; e = e.Next() {
		if fn(e.Value) {
			return e.Value, true
		}
	}

	return nil, false
}

func (sc *SList) Remove(fn func(i interface{}) bool) {
	sc.locker.Lock()
	defer sc.locker.Unlock()

	for e := sc.slist.Front(); e != nil; e = e.Next() {
		if fn(e.Value) {
			sc.slist.Remove(e)
			return
		}
	}
}

func (sc *SList) Sort(fn func(a, b interface{}) bool) {
	sc.locker.Lock()
	defer sc.locker.Unlock()

	newList := list.New()
	//用老链表进行遍历  与新链表进行表
	for e := sc.slist.Front(); e != nil; e = e.Next() {
		node := newList.Front()
		for nil != node {
			cmpVal := fn(node.Value, e.Value)
			if cmpVal {
				newList.InsertBefore(e.Value, node)
				break
			}

			node = node.Next()
		}

		//能走到这步 则表明v只能放入链表最后
		if node == nil {
			newList.PushBack(e.Value)
		}
	}

	sc.slist = newList // 新队列
}

func (sc *SList) RemoveRanges(fn func(i interface{}) bool) {
	sc.locker.Lock()
	defer sc.locker.Unlock()

	for e := sc.slist.Front(); e != nil; {
		next := e.Next()
		if fn(e.Value) {
			sc.slist.Remove(e)
		}
		e = next
	}
}

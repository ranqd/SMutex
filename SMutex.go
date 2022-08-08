package smutex

import (
	"fmt"
	"github.com/petermattis/goid"
	"sync"
	"sync/atomic"
)

type SMutex struct {
	locker sync.Mutex
	lockByGoroutine  int64 //当前拥有锁的协程ID
}

func (s *SMutex) Lock() {
	goid := goid.Get()
	//同一协程不锁定
	if atomic.LoadInt64(&s.lockByGoroutine) == goid {
		return
	}
	s.locker.Lock()
	atomic.StoreInt64(&s.lockByGoroutine, goid)
}

func (s *SMutex) UnLock() {
	goid := goid.Get()
	if s.lockByGoroutine == 0 {
		//fmt.Println("重复解锁")
		return
	}
	if s.lockByGoroutine != goid {
		fmt.Println("锁定", s.lockByGoroutine,"解锁", goid,"协程不一致")
		//return
	}
	atomic.StoreInt64(&s.lockByGoroutine, 0)
	s.locker.Unlock()
}

//TryLock 尝试锁定，成功返回true，失败返回false
func (s *SMutex) TryLock() bool {
	goid := goid.Get()
	//同一协程不锁定
	if atomic.LoadInt64(&s.lockByGoroutine) == goid {
		return true
	}
	ret := s.locker.TryLock()
	if ret {
		atomic.StoreInt64(&s.lockByGoroutine, goid)
	}
	return ret
}

type SRWMutex struct {
	locker sync.RWMutex
	rLockByGoroutine  int64 //当前拥有读锁的协程ID
	wLockByGoroutine  int64 //当前拥有写锁的协程ID
}

func (s *SRWMutex) Lock() {
	goid := goid.Get()
	//同一协程不锁定
	if atomic.LoadInt64(&s.wLockByGoroutine) == goid {
		return
	}
	s.locker.Lock()
	atomic.StoreInt64(&s.wLockByGoroutine, goid)
}

func (s *SRWMutex) UnLock() {
	goid := goid.Get()
	if s.wLockByGoroutine == 0 {
		//fmt.Println("重复解锁")
		return
	}
	if s.wLockByGoroutine != goid {
		fmt.Println("锁定", s.wLockByGoroutine,"解锁", goid,"协程不一致")
	}
	atomic.StoreInt64(&s.wLockByGoroutine, 0)
	s.locker.Unlock()
}

func (s *SRWMutex) RLock() {
	goid := goid.Get()
	//同一协程不锁定
	if atomic.LoadInt64(&s.wLockByGoroutine) == goid { //写锁在同一协程
		return
	}
	s.locker.RLock()
	atomic.StoreInt64(&s.rLockByGoroutine, goid)
}

func (s *SRWMutex) RUnLock() {
	goid := goid.Get()
	if s.rLockByGoroutine == 0 {
	//	fmt.Println("重复解锁")
		return
	}
	if s.rLockByGoroutine != goid {
		fmt.Println("锁定", s.rLockByGoroutine,"解锁", goid,"协程不一致")
		//return
	}
	atomic.StoreInt64(&s.rLockByGoroutine, 0)
	s.locker.RUnlock()
}

//TryLock 尝试锁定，成功返回true，失败返回false
func (s *SRWMutex) TryLock() bool {
	goid := goid.Get()
	//同一协程不锁定
	if atomic.LoadInt64(&s.wLockByGoroutine) == goid {
		return true
	}
	ret := s.locker.TryLock()
	if ret {
		atomic.StoreInt64(&s.wLockByGoroutine, goid)
	}
	return ret
}
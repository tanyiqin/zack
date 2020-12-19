package utils

import (
	"io"
	"os"
	"sync"
	"time"
)

// 当前时间/s
func NowTime() int64{
	t := time.Now()
	return t.Unix()
}

// 用来等待所有go协程返回
type WrapWait struct {
	sync.WaitGroup
}

func (w *WrapWait)Wrap(f func()) {
	w.Add(1)
	go func() {
		f()
		w.Done()
	}()
}

// copy文件
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.Create(dstName)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

// 信号量
type Empty interface {}
type Semaphore chan Empty
var sem = make(Semaphore, 1)

func (s Semaphore) P(n int) {
	e := new(Empty)
	for i := 0; i < n; i++ {
		s <- e
	}
}

func (s Semaphore) V(n int) {
	for i := 0; i < n; i++ {
		<- s
	}
}

/* mutexes */
func (s Semaphore) Lock() {
	s.P(1)
}
func (s Semaphore) UnLock() {
	s.V(1)
}

/* signal-wait */
func (s Semaphore) Wait (n int) {
	s.P(n)
}
func (s Semaphore) Signal() {
	s.V(1)
}
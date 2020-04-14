package factory

import (
	"fmt"
	"sync/atomic"
	"time"
)

type exitSignal struct{}

// worker 工作者角色
type worker struct {
	isBusy int32
	action func(interface{})
	params chan interface{}
}

func (w *worker) process() (quit bool) {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("worker broken, panic = %v", p)
		}
	}()
	for params := range w.params {
		if _, ok := params.(exitSignal); ok {
			return true
		}
		w.action(params)
		atomic.StoreInt32(&w.isBusy, 0)
	}
	return false
}

func (w *worker) assign(line *Line, params interface{}) bool {
	if atomic.CompareAndSwapInt32(&w.isBusy, 0, 1) {
		line.waitGroup.Add(1)
		w.action = line.action
		w.params <- params
		return true
	}
	return false
}

func (w *worker) shutdown() {
	w.params <- exitSignal{}
}

func newWorker() (w *worker) {
	w = &worker{
		params: make(chan interface{}),
	}
	go func(w *worker) {
		for {
			if w.process() {
				break
			}
			atomic.StoreInt32(&w.isBusy, 0)
		}
		// 置为繁忙状态
		atomic.StoreInt32(&w.isBusy, 1)
		// 关闭任务通道
		defer close(w.params)
		// 可能存在任务
		select {
		case params := <-w.params:
			w.action(params)
		case <-time.Tick(time.Millisecond):
		}
	}(w)
	return
}

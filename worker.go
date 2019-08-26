package factory

import (
	"fmt"
	"sync/atomic"
)

// worker 工作者角色
type worker struct {
	isBusy int32
	action func(interface{})
	params chan interface{}
}

func (w *worker) process() bool {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("worker broken, panic = %v", p)
		}
	}()
	for params := range w.params {
		if _, ok := params.(struct{}); ok {
			return true
		}
		w.action(params)
		atomic.StoreInt32(&w.isBusy, 0)
	}
	return false
}

func (w *worker) assign(action func(interface{}), params interface{}) bool {
	if atomic.CompareAndSwapInt32(&w.isBusy, 0, 1) {
		w.action = action
		w.params <- params
		return true
	}
	return false
}

func (w *worker) shutdown() {
	w.params <- struct{}{}
}

func newWorker() (w *worker) {
	w = &worker{
		params: make(chan interface{}),
	}
	go func(w *worker) {
		for w.process() == false {
			atomic.StoreInt32(&w.isBusy, 0)
		}
		// 置为繁忙状态
		atomic.StoreInt32(&w.isBusy, 1)
		// 可能后入任务
		select {
		case params := <-w.params:
			w.action(params)
		default:
		}
	}(w)
	return
}

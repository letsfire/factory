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

func (w *worker) process() (quit bool) {
	defer func(quit *bool) {
		if p := recover(); p != nil {
			*quit = false
			fmt.Printf("worker broken, panic = %v", p)
		}
	}(&quit)
	for params := range w.params {
		w.action(params)
		atomic.StoreInt32(&w.isBusy, 0)
	}
	return true
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
	// 置为繁忙状态
	atomic.StoreInt32(&w.isBusy, 1)
	// 可能存在任务
	select {
	case params := <-w.params:
		w.action(params)
	default:
	}
	// 关闭任务通道
	close(w.params)
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
	}(w)
	return
}

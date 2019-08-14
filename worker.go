package factory

import (
	"log"
)

// worker 工作者角色
type worker struct {
	line *Line
	args chan interface{}
}

func (w *worker) process() bool {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("worker broken, panic = %#v", p)
		}
	}()
	for args := range w.args {
		if _, ok := args.(struct{}); ok {
			return true
		}
		w.line.action(args)
	}
	return false
}

func (w *worker) assign(l *Line, args interface{}) {
	w.line = l
	w.args <- args
}

func (w *worker) shutdown() {
	w.args <- struct{}{}
}

func newWorker() (w *worker) {
	w = &worker{
		args: make(chan interface{}),
	}
	go func() {
		for {
			if quit := w.process(); quit {
				break
			}
		}
	}()
	return w
}

package factory

import (
	"log"
)

// worker 工作者角色
type worker struct {
	master *Master
}

func (w worker) process() {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("worker broken from panic, %#v", p)
		}
	}()
	var t task
	for {
		select {
		case t = <-w.master.topChan:
		default:
			select {
			case t = <-w.master.topChan:
			case t = <-w.master.middleChan:
			default:
				select {
				case t = <-w.master.topChan:
				case t = <-w.master.middleChan:
				case t = <-w.master.bottomChan:
				}
			}
		}
		t.action(t.params...)
	}
}

func newWorker(master *Master) (w worker) {
	w = worker{
		master: master,
	}
	go func() {
		for {
			w.process()
		}
	}()
	return w
}

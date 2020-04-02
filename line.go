package factory

import (
	"math"
	"sync"
)

// Line 工作流水线
type Line struct {
	sync.WaitGroup
	name   string
	master *Master
	action func(interface{})
}

func (l *Line) Submit(args interface{}) {
	if l.master.getWorker().assign(l.action, args) {
		return
	}
	tryNum := int(float64(l.master.Running()) * 0.5)
	for n := 0; n < tryNum; n++ {
		if l.master.getWorker().assign(l.action, args) {
			return
		}
	}
	addSize := math.Max(2, float64(l.master.Running())*0.25)
	l.master.AdjustSize(l.master.Running() + int(addSize))
	l.Submit(args)
}

func (l *Line) Execute(args interface{}) { l.action(args) }

func (l *Line) SetPanicHandler(handler func(interface{})) {
	oldAction := l.action
	l.action = func(i interface{}) {
		defer func() {
			if p := recover(); p != nil {
				handler(p)
			}
		}()
		oldAction(i)
	}
}

func NewLine(master *Master, name string, action func(interface{})) *Line {
	l := Line{master: master, name: name}
	l.action = func(i interface{}) {
		l.Add(1)
		action(i)
		l.Done()
	}
	return &l
}

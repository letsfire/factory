package factory

import (
	"math"
	"sync"
)

// Line 工作流水线
type Line struct {
	name      string
	master    *Master
	action    func(interface{})
	waitGroup *sync.WaitGroup
}

func (l *Line) Wait() {
	l.waitGroup.Wait()
}

func (l *Line) Submit(args interface{}) {
	if l.master.getWorker().assign(l.action, args) {
		l.waitGroup.Add(1)
		return
	}
	tryNum := int(float64(l.master.Running()) * 0.5)
	for n := 0; n < tryNum; n++ {
		if l.master.getWorker().assign(l.action, args) {
			l.waitGroup.Add(1)
			return
		}
	}
	addSize := math.Max(2, float64(l.master.Running())*0.25)
	l.master.AdjustSize(l.master.Running() + int(addSize))
	l.Submit(args)
}

func (l *Line) Execute(args interface{}) {
	l.waitGroup.Add(1)
	l.action(args)
}

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
	wg := new(sync.WaitGroup)
	return &Line{
		name:   name,
		master: master,
		action: func(i interface{}) {
			defer wg.Done()
			action(i)
		},
		waitGroup: wg,
	}
}

package factory

import (
	"math"
)

// Line 工作流水线
type Line struct {
	name   string
	master *Master
	action func(interface{})
}

func (l *Line) Submit(args interface{}) {
	if l.master.getWorker().assign(l.action, args) {
		return
	}
	tryNum := int(float64(l.master.ingNum) * 0.5)
	for n := 0; n < tryNum; n++ {
		if l.master.getWorker().assign(l.action, args) {
			return
		}
	}
	addSize := math.Max(2, float64(l.master.ingNum)*0.25)
	l.master.AdjustSize(int(l.master.ingNum) + int(addSize))
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

func NewLine(master *Master, name string, action func(interface{})) (l *Line) {
	return &Line{master: master, name: name, action: action}
}

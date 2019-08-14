package factory

// Line 工作流水线
type Line struct {
	name   string
	master *Master
	action func(interface{})
}

func (l *Line) Submit(args interface{}) {
	l.master.getWorker().assign(l, args)
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

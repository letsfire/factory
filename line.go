package factory

// Line 工作流水线
type Line struct {
	name   string
	level  Level
	master *Master
	action func(...interface{})
}

func (l *Line) Submit(args ...interface{}) {
	switch l.level {
	case Top:
		l.SubmitTop(args...)
	case Middle:
		l.SubmitMiddle(args...)
	default:
		l.SubmitBottom(args...)
	}
}

func (l *Line) SubmitTop(args ...interface{}) {
	l.master.topChan <- task{line: l, args: args}
}

func (l *Line) SubmitMiddle(args ...interface{}) {
	l.master.middleChan <- task{line: l, args: args}
}

func (l *Line) SubmitBottom(args ...interface{}) {
	l.master.bottomChan <- task{line: l, args: args}
}

func (l *Line) SetLevel(lv Level)           { l.level = lv }
func (l *Line) Execute(args ...interface{}) { l.action(args...) }

func NewLine(master *Master, name string, action func(...interface{})) (l *Line) {
	return &Line{
		master: master, name: name,
		action: action, level: Bottom,
	}
}

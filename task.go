package factory

// task 流水线任务
type task struct {
	line *Line
	args []interface{}
}

func (t *task) execute() {
	t.line.action(t.args...)
}

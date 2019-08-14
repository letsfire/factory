package factory

// task 流水线任务
type task struct {
	action func(...interface{})
	params []interface{}
}

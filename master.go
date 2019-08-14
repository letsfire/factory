package factory

// Master 管理者角色
type Master struct {
	workers    []worker
	topChan    chan task
	middleChan chan task
	bottomChan chan task
}

func NewMaster(size int) (m *Master) {
	m = &Master{
		workers:    make([]worker, size),
		topChan:    make(chan task, 1),
		middleChan: make(chan task, 1),
		bottomChan: make(chan task, 1),
	}
	for i := 0; i < size; i++ {
		m.workers[i] = newWorker(m)
	}
	return
}

func (m *Master) AddLine(name string, action func(...interface{})) *Line {
	return NewLine(m, name, action)
}

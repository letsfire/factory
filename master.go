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
		topChan:    make(chan task, size>>2),
		middleChan: make(chan task, size>>2),
		bottomChan: make(chan task, size>>1),
	}
	for i := 0; i < size; i++ {
		m.workers[i] = newWorker(m)
	}
	return
}

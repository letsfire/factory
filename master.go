package factory

import (
	"sync"
	"sync/atomic"
)

// Master 管理者角色
type Master struct {
	sync.RWMutex
	size    int32
	cursor  int32
	workers []*worker
}

func NewMaster(size int) (m *Master) {
	m = &Master{
		size:    int32(size),
		workers: make([]*worker, size),
	}
	for i := 0; i < size; i++ {
		m.workers[i] = newWorker()
	}
	return
}

func (m *Master) AdjustSize(newSize int) {
	diff := newSize - len(m.workers)
	if diff == 0 {
		return
	}
	m.Lock()
	defer m.Unlock()
	if diff > 0 {
		// 扩容
		for i := 0; i < diff; i++ {
			m.workers = append(m.workers, newWorker())
		}
		m.size = int32(newSize)
	} else {
		// 缩容
		m.size = int32(newSize)
		dws := m.workers[0:-diff]
		m.workers = m.workers[-diff:]
		for _, w := range dws {
			w.shutdown()
		}
	}
}

func (m *Master) AddLine(name string, action func(interface{})) *Line {
	return NewLine(m, name, action)
}

func (m *Master) getWorker() *worker {
	m.RLock()
	atomic.CompareAndSwapInt32(&m.cursor, m.size, 0)
	idx := atomic.AddInt32(&m.cursor, 1)
	w := m.workers[idx-1]
	m.RUnlock()
	return w
}

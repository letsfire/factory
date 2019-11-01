package factory

import (
	"sync"
	"sync/atomic"
)

// Master 管理者角色
type Master struct {
	sync.Mutex
	maxNum  int64
	ingNum  int64
	cursor  int64
	workers []*worker
}

func NewMaster(maxNum, initNum int) (m *Master) {
	m = &Master{
		maxNum:  int64(maxNum),
		ingNum:  int64(initNum),
		workers: make([]*worker, maxNum),
	}
	for i := 0; i < initNum; i++ {
		m.workers[i] = newWorker()
	}
	return
}

func (m *Master) AddLine(name string, action func(interface{})) *Line {
	return NewLine(m, name, action)
}

func (m *Master) AdjustSize(newSize int) {
	if int64(newSize) > m.maxNum {
		newSize = int(m.maxNum)
	}

	m.Lock()
	defer m.Unlock()

	if diff := newSize - int(m.ingNum); diff > 0 {
		for i := 0; i < diff; i++ {
			m.workers[m.ingNum] = newWorker()
			atomic.AddInt64(&m.ingNum, 1)
		}
	} else if diff < 0 {
		atomic.StoreInt64(&m.ingNum, int64(newSize))
		if cursor := atomic.LoadInt64(&m.cursor); cursor > int64(newSize) {
			atomic.StoreInt64(&m.cursor, int64(newSize))
		}
		for idx, w := range m.workers[newSize:] {
			if w == nil {
				break
			}
			w.shutdown()
			m.workers[idx] = nil
		}
	}
}

func (m *Master) Running() int64 {
	return atomic.LoadInt64(&m.ingNum)
}

func (m *Master) Shutdown() {
	m.AdjustSize(0) // 关闭所有worker
}

func (m *Master) getWorker() *worker {
	atomic.CompareAndSwapInt64(&m.cursor, m.ingNum, 0)
	idx := atomic.AddInt64(&m.cursor, 1)
	w := m.workers[idx-1]
	return w
}

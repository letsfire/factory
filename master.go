package factory

import (
	"sync"
	"sync/atomic"
)

// Master 管理者角色
type Master struct {
	sync.Mutex
	maxNum  int64
	cursor  int64
	idxFlag int32
	tryFlag int32
	synCond *sync.Cond
	workers [2]map[int64]*worker
}

func NewMaster(maxNum, initNum int) *Master {
	m := new(Master)
	m.maxNum = int64(maxNum)
	m.workers[0] = make(map[int64]*worker, initNum)
	for i := 0; i < initNum; i++ {
		m.workers[0][int64(i)] = newWorker()
	}
	m.workers[1] = m.workers[0]
	m.synCond = sync.NewCond(&sync.Mutex{})
	return m
}

func (m *Master) AddLine(name string, action func(interface{})) *Line {
	return NewLine(m, name, action)
}

func (m *Master) AdjustSize(newSize int) {
	if newSize > int(m.maxNum) {
		newSize = int(m.maxNum)
	}

	m.Lock()
	defer m.Unlock()

	workers := m.usingGroup()
	diff := newSize - len(workers)

	if diff > 0 {
		for i := len(workers); i < newSize; i++ {
			workers[int64(i)] = newWorker()
		}
	} else if diff < 0 {
		for i := len(workers); i > newSize; i-- {
			idx := int64(i) - 1
			workers[idx].shutdown()
			delete(workers, idx)
		}
	}

	idxFlag := 2 + ^m.idxFlag
	m.workers[idxFlag] = workers
	atomic.StoreInt32(&m.idxFlag, idxFlag)
	m.workers[2+^m.idxFlag] = m.workers[idxFlag]
}

func (m *Master) Running() int {
	return len(m.usingGroup())
}

func (m *Master) Shutdown() {
	m.AdjustSize(0) // 关闭所有worker
}

func (m *Master) usingGroup() map[int64]*worker {
	return m.workers[atomic.LoadInt32(&m.idxFlag)]
}

func (m *Master) getWorker() *worker {
	idx := atomic.AddInt64(&m.cursor, 1)
	if w, ok := m.usingGroup()[idx-1]; ok && w != nil {
		return w
	}
	if atomic.CompareAndSwapInt32(&m.tryFlag, 0, 1) {
		atomic.StoreInt64(&m.cursor, 0)
		atomic.StoreInt32(&m.tryFlag, 0)
		m.synCond.Broadcast()
	} else {
		m.synCond.L.Lock()
		m.synCond.Wait()
		m.synCond.L.Unlock()
	}
	return m.getWorker()
}

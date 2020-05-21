package factory

import (
	"sync"
	"sync/atomic"

	"github.com/letsfire/utils"
)

// Master 管理者角色
type Master struct {
	sync.Mutex
	maxNum     int
	ingNum     int
	cursor     int64
	workers    sync.Map
	resetGuard *utils.Guard
}

func NewMaster(maxNum, initNum int) *Master {
	m := new(Master)
	m.maxNum = maxNum
	m.ingNum = initNum
	for i := 0; i < initNum; i++ {
		m.workers.Store(i, newWorker())
	}
	m.resetGuard = utils.NewGuard()
	return m
}

func (m *Master) Resize(maxNum int) {
	m.maxNum = maxNum
	if m.ingNum > maxNum {
		m.AdjustSize(maxNum)
	}
}

func (m *Master) AddLine(action func(interface{})) *Line {
	return NewLine(m, action)
}

func (m *Master) AdjustSize(newSize int) {
	if newSize > m.maxNum {
		newSize = m.maxNum
	}
	m.Lock()
	defer m.Unlock()
	if newSize > m.ingNum {
		for i := m.ingNum; i < newSize; i++ {
			m.workers.Store(i, newWorker())
		}
	} else if newSize < m.ingNum {
		for i := m.ingNum; i > newSize; i-- {
			idx := i - 1
			if v, ok := m.workers.Load(idx); ok {
				m.workers.Delete(idx)
				v.(*worker).shutdown()
			}
		}
	}
	m.ingNum = newSize
}

func (m *Master) Running() int {
	m.Lock()
	defer m.Unlock()
	return m.ingNum
}

func (m *Master) Shutdown() {
	m.AdjustSize(0) // 关闭所有worker
}

func (m *Master) getWorker() *worker {
	idx := int(atomic.AddInt64(&m.cursor, 1)) - 1
	if w, ok := m.workers.Load(idx); ok && w != nil {
		return w.(*worker)
	} else if m.ingNum == 0 {
		panic("factory: the master has been shutdown")
	}
	m.resetGuard.Run("get-worker", func() (i interface{}, e error) {
		atomic.StoreInt64(&m.cursor, 0)
		return nil, nil
	})
	return m.getWorker()
}

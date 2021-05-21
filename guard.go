package factory

import "sync"

// guard 防止并发
type guard struct {
	locker    sync.Mutex
	callerMap map[string]*caller
}

func newGuard() *guard {
	return &guard{
		callerMap: make(map[string]*caller),
	}
}

func (g *guard) run(key string, callback func() (interface{}, error)) (interface{}, error) {
	g.locker.Lock()
	c, ok := g.callerMap[key]
	if ok {
		g.locker.Unlock()
		c.waiter.Wait()
		return c.result()
	} else {
		c = newCall()
		g.callerMap[key] = c
		g.locker.Unlock()
	}
	c.run(callback)
	g.locker.Lock()
	delete(g.callerMap, key)
	g.locker.Unlock()
	return c.result()
}

type caller struct {
	value  interface{}
	error  error
	waiter sync.WaitGroup
}

func newCall() *caller {
	c := new(caller)
	c.waiter.Add(1)
	return c
}

func (c *caller) run(fn func() (interface{}, error)) {
	c.value, c.error = fn()
	c.waiter.Done()
}

func (c *caller) result() (interface{}, error) {
	return c.value, c.error
}

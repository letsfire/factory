package factory

import (
	"sync"
	"sync/atomic"
	"time"
)

var swg = new(sync.WaitGroup)
var mst = NewMaster(50000, 5000)

var counter int64
var runTimes = 50000

var lineFunc1 = func(v interface{}) {
	atomic.AddInt64(&counter, int64(v.(int)))
	time.Sleep(time.Microsecond)
}

var lineFunc2 = func(v interface{}) {
	atomic.AddInt64(&counter, -int64(v.(int)))
	time.Sleep(time.Microsecond)
}

var line1 = mst.AddLine("test.line.1", lineFunc1)
var line2 = mst.AddLine("test.line.2", lineFunc2)

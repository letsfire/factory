package factory

import (
	"sync/atomic"
)

var mst = NewMaster(50000, 5000)

var counter int64

var lineFunc1 = func(v interface{}) {
	atomic.AddInt64(&counter, int64(v.(int)))
}

var lineFunc2 = func(v interface{}) {
	atomic.AddInt64(&counter, -int64(v.(int)))
}

var line1 = mst.AddLine(lineFunc1)
var line2 = mst.AddLine(lineFunc2)

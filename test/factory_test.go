package test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/letsfire/factory"
)

var swg = new(sync.WaitGroup)
var mst = factory.NewMaster(100000)

var counter int64

var lineFunc1 = func(args ...interface{}) {
	atomic.AddInt64(&counter, int64(args[0].(int)))
	time.Sleep(time.Millisecond * 10)
	swg.Done()
}

var lineFunc2 = func(args ...interface{}) {
	atomic.AddInt64(&counter, -int64(args[0].(int)))
	time.Sleep(time.Millisecond * 5)
	swg.Done()
}

var line1 = mst.AddLine("test.line.1", lineFunc1)
var line2 = mst.AddLine("test.line.2", lineFunc2)

func BenchmarkWithFactory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		swg.Add(2)
		line1.SubmitTop(1)
		line2.SubmitMiddle(1)
	}
	swg.Wait()
	if counter != 0 {
		b.Errorf("unexpect  result, expect = 0, but = %d", counter)
	}
}

func BenchmarkWithOutFactory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		swg.Add(2)
		go lineFunc1(1)
		go lineFunc2(1)
	}
	swg.Wait()
	if counter != 0 {
		b.Errorf("unexpect  result, expect = 0, but = %d", counter)
	}
}

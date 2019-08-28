package issue

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/letsfire/factory"
)

type hh struct {
	mux *sync.Mutex
	a   int
}

//go test -v -test.run TestNewMaster
func TestNewMaster(t *testing.T) {
	master := factory.NewMaster(1000, 1000)
	// 新建第一条工作流水线
	var line1 = master.AddLine("demo.line.1", func(args interface{}) {
		h := args.(hh)
		h.mux.Lock()
		h.a++
		if h.a%10000 == 0 {
			fmt.Println(h.a)
		}
		h.mux.Unlock()
		time.Sleep(time.Millisecond * 10)
	})

	// 根据业务场景将参数提交
	mux := &sync.Mutex{}
	h := hh{mux: mux}
	for i := 0; i < 1000000; i++ {
		h.a = i
		line1.Submit(h)
	}

	// 协程池数量可动态调整
	master.Running() // 正在运行的协程工人数量
	//master.AdjustSize(100)      // 指定数量进行扩容或缩容
	master.Shutdown() // 等于 master.AdjustSize(0)
	fmt.Println(runtime.NumGoroutine())
	time.Sleep(time.Second * 5)
	fmt.Println(runtime.NumGoroutine())
}

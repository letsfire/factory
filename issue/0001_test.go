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
	master := factory.NewMaster(10000, 10)
	// 新建第一条工作流水线
	var line1 = master.AddLine(func(args interface{}) {
		h := args.(hh)
		h.mux.Lock()
		h.a++
		if h.a%10000 == 0 {
			// fmt.Println(h.a)
		}
		h.mux.Unlock()
		time.Sleep(time.Second)
	})

	// 根据业务场景将参数提交
	mux := &sync.Mutex{}
	h := hh{mux: mux}
	for i := 0; i < 10000; i++ {
		h.a = i
		line1.Submit(h)
	}

	// 协程池数量可动态调整
	fmt.Println(master.Running()) // 正在运行的协程工人数量

	master.AdjustSize(1000)       // 指定数量进行扩容或缩容
	fmt.Println(master.Running()) // 正在运行的协程工人数量

	master.Shutdown()             // 等于 master.AdjustSize(0)
	fmt.Println(master.Running()) // 正在运行的协程工人数量

	fmt.Println(runtime.NumGoroutine()) // 通常会大于 2
	time.Sleep(time.Millisecond * 5)    // 极短暂休眠
	fmt.Println(runtime.NumGoroutine()) // 必须会等于 2
}

package issue

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/letsfire/factory"
)

func TestMaster(t *testing.T) {
	wg := sync.WaitGroup{}
	wg2 := sync.WaitGroup{}
	w := factory.NewMaster(10, 2)
	ctx, cel := context.WithCancel(context.Background())
	li := w.AddLine(func(e interface{}) {
		defer wg.Done()
	})
	for i := 0; i < 100; i++ {
		go func() {
			wg2.Add(1)
			defer wg2.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					wg.Add(1)
					li.Submit(0)
				}
			}
		}()
	}
	time.Sleep(time.Second * 5)
	fmt.Println("adjustsize")
	w.AdjustSize(6)
	time.Sleep(time.Second * 2)
	w.AdjustSize(1)
	time.Sleep(time.Second * 2)
	cel() //关闭发协程
	w.AdjustSize(2)
	time.Sleep(time.Second * 2)
	wg2.Wait() //确认发协程是否能关闭，多半会卡在这里
	fmt.Println("down")
	w.Shutdown()
	fmt.Println("Wait")
	wg.Wait()
	fmt.Println("End")
}

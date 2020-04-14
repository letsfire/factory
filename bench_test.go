package factory

import (
	"testing"
	"time"
)

func BenchmarkWithGoroutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		go line1.Execute(1)
		go line2.Execute(1)
	}
	b.StopTimer()
	time.Sleep(time.Second)
	b.StartTimer()
	line1.Wait()
	line2.Wait()
	if counter != 0 {
		b.Errorf("unexpect  result, expect = 0, but = %d", counter)
	}
}

func BenchmarkWithFactory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		line1.Submit(1)
		line2.Submit(1)
	}
	line1.Wait()
	line2.Wait()
	if counter != 0 {
		b.Errorf("unexpect  result, expect = 0, but = %d", counter)
	}
}

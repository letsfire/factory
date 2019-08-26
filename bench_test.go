package factory

import (
	"testing"
)

func BenchmarkWithGoroutine(b *testing.B) {
	swg.Add(2 * b.N)
	for i := 0; i < b.N; i++ {
		go line1.Execute(1)
		go line2.Execute(1)
	}
	swg.Wait()
	if counter != 0 {
		b.Errorf("unexpect  result, expect = 0, but = %d", counter)
	}
}

func BenchmarkWithFactory(b *testing.B) {
	swg.Add(2 * b.N)
	for i := 0; i < b.N; i++ {
		line1.Submit(1)
		line2.Submit(1)
	}
	swg.Wait()
	if counter != 0 {
		b.Errorf("unexpect  result, expect = 0, but = %d", counter)
	}
}

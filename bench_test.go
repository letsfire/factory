package factory

import (
	"testing"
)

func BenchmarkWithGoroutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		swg.Add(2 * runTimes)
		for j := 0; j < runTimes; j++ {
			go line1.Execute(1)
			go line2.Execute(1)
		}
	}
	swg.Wait()
	if counter != 0 {
		b.Errorf("unexpect  result, expect = 0, but = %d", counter)
	}
}

func BenchmarkWithFactory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		swg.Add(2 * runTimes)
		for j := 0; j < runTimes; j++ {
			line1.Submit(1)
			line2.Submit(1)
		}
	}
	swg.Wait()
	if counter != 0 {
		b.Errorf("unexpect  result, expect = 0, but = %d", counter)
	}
}

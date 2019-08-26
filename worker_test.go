package factory

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {
	w := newWorker()

	// panic test

	/*
		swg.Add(1)
		assert.True(t, w.assign(lineFunc1, "x"))
		swg.Wait()
		time.Sleep(time.Millisecond * 100)
		assert.EqualValues(t, 0, w.isBusy)
	*/

	swg.Add(1)
	assert.True(t, w.assign(lineFunc1, 1))
	swg.Wait()
	assert.EqualValues(t, 1, counter)

	swg.Add(1)
	assert.True(t, w.assign(lineFunc2, 1))
	swg.Wait()
	assert.EqualValues(t, 0, counter)

	w.shutdown()
	time.Sleep(time.Millisecond * 100)
	assert.EqualValues(t, 1, atomic.LoadInt32(&w.isBusy))
}

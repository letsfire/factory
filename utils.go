package factory

import (
	"math"
	"runtime"
)

// handlePanic 使用回调函数处理panic
func handlePanic(hs ...func(interface{})) {
	if r := recover(); r != nil {
		for _, handler := range hs {
			handler(r)
		}
	}
}

// stackTrace 获取堆栈踪迹
func stackTrace(size int) []byte {
	if size <= 0 {
		size = math.MaxUint16
	}
	stacktrace := make([]byte, size)
	return stacktrace[:runtime.Stack(stacktrace, false)]
}

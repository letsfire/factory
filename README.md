## Factory
Go语言的协程池 , 节省内存 , 减少GC压力

## 安装
`go get github.com/letsfire/factory`

## 用法
```go
// 新建一个协程池,指定协程数量10000
var master = factory.NewMaster(10000)

// 新建第一条工作流水线
var line1 = factory.NewLine(master, "demo.line.1", func(args ...interface{}) {

	// TODO 处理您的业务逻辑
	// fmt.Println(args...)

})

// 新建第二条工作流水线
var line2 = factory.NewLine(master, "demo.line.2", func(args ...interface{}) {

	// TODO 处理您的业务逻辑
	// fmt.Println(args...)

})

// 根据业务场景将参数提交
for i := 0; i < 100000; i++ {
	line1.SubmitTop(1, 2, 3)		// 最高优先级
	line1.SubmitMiddle(1, 2, 3)		// 中等优先级
	line1.SubmitBottom(1, 2, 3)		// 最低优先级
}

// 默认流水线的优先级是Bottom , 您可以动态调整
line2.SetLevel(Top)
line2.Submit(1, 2)	// 等同于 line2.SubmitTop(1, 2)

// SubmitTop, SubmitMiddle, SubmitBottom 三个方法可以覆盖流水线的默认优先级
```
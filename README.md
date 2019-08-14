## Factory
Go语言的协程池 , 节省内存 , 减少GC压力

## 安装
`go get github.com/letsfire/factory`

## 用法
```go
// 新建一个协程池,指定协程数量20000
var master = factory.NewMaster(20000)

// 新建第一条工作流水线
var line1 = master.AddLine("demo.line.1", func(args interface{}) {

	// TODO 处理您的业务逻辑
	// fmt.Println(args...)

})

// 新建第二条工作流水线
var line2 = master.AddLine("demo.line.2", func(args interface{}) {

	// TODO 处理您的业务逻辑
	// fmt.Println(args...)

})

// 根据业务场景将参数提交
for i := 0; i < 100000; i++ {
	line1.Submit(i)
}

for j := 0; j < 100000; j++ {
	line2.Submit(j)
}

// 协程池数量可动态调整
master.AdjustSize(10000)    // 缩容
master.AdjustSize(30000)    // 扩容
```
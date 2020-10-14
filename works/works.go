package works

// IWork 对外处理
type IWork interface {
	// 处理函数
	Process() (IWork, error)
}

// IWorkChan 队列类型
type IWorkChan chan IWork

// NewChannel 创建任务列队
func NewChannel(num int) IWorkChan {
	return make(IWorkChan, num)
}

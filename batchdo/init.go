package batchdo

import "time"

const (
	// DefaultMaxTimeInv 最大提交时间间隔
	defaultMaxTimeInv = 10 * time.Second
	// DefaultMaxCount 最大提交计数
	defaultMaxCount = 1024
)

// IBatchdo .
type IBatchdo interface {
	Add(v interface{}) IBatchdo

	DoCondition(count int32, timeinv time.Duration) IBatchdo

	DoCallback(docb func(dos []interface{}) error) IBatchdo

	Erorr() (errs <-chan error)
}

// New 批量执行
func New() IBatchdo {
	b := &batch{
		maxCount:   defaultMaxCount,
		maxTimeInv: defaultMaxTimeInv,
		chdos:      make(chan []interface{}, 32),
	}
	go b.run()
	return b
}

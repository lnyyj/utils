package batchdo

import "time"

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

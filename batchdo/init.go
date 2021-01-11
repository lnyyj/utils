package batchdo

import "time"

type IBatchdo interface {
	Add(v interface{}) IBatchdo

	DoCondition(count int32, timeinv time.Duration) IBatchdo

	DoCallback(docb func(dos []interface{}) error) IBatchdo

	Erorr() (errs <-chan error)
}

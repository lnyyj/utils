package batchdo

import (
	"sync"
	"time"
)

const (
	// DefaultMaxTimeInv 最大提交时间间隔
	defaultMaxTimeInv = 10 * time.Second
	// DefaultMaxCount 最大提交计数
	defaultMaxCount = 1024
)

type batch struct {
	dos        []interface{}
	doCallback func(dos []interface{}) error

	maxCount   int32
	maxTimeInv time.Duration

	errors chan error
	chdos  chan []interface{}

	sync.Mutex
}

// NewBatchDo 批量执行
func NewBatchDo() IBatchdo {
	b := &batch{
		maxCount:   defaultMaxCount,
		maxTimeInv: defaultMaxTimeInv,
		chdos:      make(chan []interface{}, 32),
	}
	go b.run()
	return b
}

func (b *batch) DoCondition(count int32, timeinv time.Duration) IBatchdo {
	b.maxCount = count
	b.maxTimeInv = timeinv
	return b
}
func (b *batch) DoCallback(docb func(dos []interface{}) error) IBatchdo {
	b.doCallback = docb
	return b
}
func (b *batch) Erorr() (errs <-chan error) {
	b.errors = make(chan error)
	return b.errors
}

func (b *batch) Add(v interface{}) IBatchdo {
	if count := len(b.dos); int32(count) >= b.maxCount {
		b.addChdos()
	}
	b.dos = append(b.dos, v)
	return b
}

func (b *batch) addChdos() {
	b.Lock()
	defer b.Unlock()
	if l := len(b.dos); l > 0 {
		b.chdos <- b.dos
		b.dos = make([]interface{}, 0)
	}
}

func (b *batch) run() {
	for {
		select {
		case <-time.After(b.maxTimeInv):
			b.addChdos()
		case dos := <-b.chdos:
			if b.doCallback != nil {
				if err := b.doCallback(dos); err != nil && b.errors != nil {
					b.errors <- err
				}
			}
		}
	}
}
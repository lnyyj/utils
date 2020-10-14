package thriftlink

import (
	"container/list"
	"sync"
	"time"
)

type TPool struct {
	// 创建对象并连接函数
	NewDial func() (interface{}, error)

	// 最大空闲对象
	MaxIdle int

	// 最大激活对象
	MaxActive int

	// 空闲对象检测超时
	IdleTimeout time.Duration

	// 信号量
	mu     sync.Mutex
	cond   *sync.Cond
	closed bool
	active int

	// 空闲客户端对象
	idle list.List
}

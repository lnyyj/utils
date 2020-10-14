package works

import (
	"time"
)

// Dispatcher 调度
type Dispatcher struct {
	// 工作池数量
	MaxWorkers int

	// 工作池
	Workers chan IWorkChan

	// 任务进入队列 必须有
	InQueue IWorkChan

	// 任务出去队列 可以为nil
	OutQueue IWorkChan

	Quit chan bool
}

// NewDispatcher 创建一个调度员
func NewDispatcher(maxWorkers int, inq, outq IWorkChan) *Dispatcher {
	workers := make(chan IWorkChan, maxWorkers)
	return &Dispatcher{Workers: workers, MaxWorkers: maxWorkers, InQueue: inq, OutQueue: outq}
}

// Run 启动调度并派遣工作
func (d *Dispatcher) Run() *Dispatcher {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := newWorker(d.Workers, d.OutQueue)
		worker.Start()
	}

	go d.dispatch()

	return d
}

// Close 关闭
func (d *Dispatcher) Close() {
	go func() {
		d.stop()
		for {
			select {
			case w := <-d.Workers:
				close(w)
			case <-time.After(1 * time.Second):
				goto _ret
			}
		}
	_ret:
		close(d.Workers)
	}()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.InQueue:
			go func(job IWork) {
				work := <-d.Workers
				work <- job
			}(job)

		case <-d.Quit:
			return
		}
	}
}

func (d *Dispatcher) stop() {
	go func() {
		d.Quit <- true
	}()
}

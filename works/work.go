package works

type worker struct {
	Workers chan IWorkChan
	Worker  IWorkChan
	Quit    chan bool

	OutQueue IWorkChan
}

// NewWork 创建一个
func newWorker(works chan IWorkChan, outq IWorkChan) *worker {
	return &worker{
		Workers:  works,
		OutQueue: outq,
		Worker:   make(chan IWork),
		Quit:     make(chan bool),
	}
}

func (w *worker) Start() {
	go func() {
		for {
			w.Workers <- w.Worker
			select {
			case job := <-w.Worker:
				if job == nil {
					return
				}
				out, _ := job.Process()
				if w.OutQueue != nil && out != nil {
					w.OutQueue <- out
				}

			case <-w.Quit:
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w *worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}

package work

import "sync"

//Worker必须满足接口类型
//才能使用工作池

type Worker interface {
	Task()
}

//Pool提供了一个goroutine池，这个池可以完成任何
//已提交的worker任务

type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

//New创建一个新工作池
func New(maxGoroutines int) *Pool {
	p := Pool{
		work: make(chan Worker),
	}

	p.wg.Add(maxGoroutines)

	for i := 0; i < maxGoroutines; i++ {
		go func() {
			defer p.wg.Done()

			//无缓冲通道，这里会一直死循环，直到通道被关闭
			for w := range p.work {
				w.Task()
			}

		}()
	}
	return &p
}

//Run提交工作到工作池
func (p *Pool) Run(w Worker) {
	p.work <- w
}

//shutdown等待所有goroutine停止工作
func (p *Pool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}

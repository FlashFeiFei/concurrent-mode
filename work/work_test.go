package work

import (
	"log"
	"sync"
	"testing"
	"time"
)

var names = []string{
	"steve",
	"bob",
	"mary",
	"therese",
	"jason",
}

//namePrinter使用特定方式打印名字
type namePrinter struct {
	name string
}

// Task实现了worker接口
func (m *namePrinter) Task() {
	log.Println(m.name)
	time.Sleep(time.Second)
}

func TestWork(t *testing.T) {
	//使用两个goroutine来创建工作池
	p := New(2)
	var wg sync.WaitGroup
	wg.Add(100 * len(names))

	for i := 0; i < 100; i++ {
		for _, name := range names {
			//创建一个namePrinter并提供
			//指定的名字
			np := namePrinter{
				name: name,
			}

			go func() {
				defer wg.Done()

				//将任务提交执行，当run返回时
				//我们就直到任务已经处理完成
				p.Run(&np)
			}()
		}
	}

	wg.Wait()

	//让工作池停止工作，等待所有现有的
	//工作完成
	p.Shutdown()
}

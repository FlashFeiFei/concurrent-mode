package runner

import (
	"log"
	"os"
	"testing"
	"time"
)

// timeout 规定了必须在多少秒内处理完成
const timeout = 3 * time.Second

func TestRunner(t *testing.T) {
	log.Println("Starting work.")

	// 为本次执行分配超时时间
	r := New(timeout)

	//加入要执行的任务
	r.Add(createTask(), createTask(), createTask())

	//执行任务并处理结果
	if err := r.Start(); err != nil {
		switch err {
		case ErrTimeout:
			log.Println("任务超时")
			os.Exit(1)
		case ErrInterrupt:
			log.Println("接收到中断信号，停止执行任务")
			os.Exit(1)
		}
	}
	log.Println("正常完成")
}

//createTask返回一个根据id
//休眠指定描述的示例任务
func createTask() func(int) {
	return func(id int) {
		log.Printf("Processor - Task #%d.", id)
		time.Sleep(time.Duration(id) * time.Second)
	}
}

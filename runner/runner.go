package runner

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

//runner 包管理处理任务的运行和声明周期

// Runner在给定的超时时间内执行一组任务
// 并且在操作系统发送中断信号时结束这些任务
type Runner struct {
	// interrupt 通道报告从操作系统
	// 发送的信号
	interrupt chan os.Signal

	// complete 通道报告处理任务已经完成
	complete chan error

	// timeout 报告处理任务已经超时
	//定义一个通道，这个通道只能是接收端
	timeout <-chan time.Time

	// tasks 持有一组以索引顺序依次执行的
	// 函数
	tasks []func(int)
}

// ErrTimeout会在任务执行超时时返回
var ErrTimeout = errors.New("received timeout")

// ErrInterrupt 会在接收到操作系统的事件时返回
var ErrInterrupt = errors.New("received interrupt")

//New 返回一个新的准备使用的Runner
func New(d time.Duration) *Runner {
	return &Runner{
		interrupt: make(chan os.Signal, 1),
		complete:  make(chan error),
		timeout:   time.After(d),
	}
}

//Add将一个任务附加到Runner上。这个任务是一个
//接受一个int类型的ID作为参数的函数
func (r *Runner) Add(tasks ...func(int)) {
	r.tasks = append(r.tasks, tasks...)
}

//gotInterrupt验证是否接收到了中断信号
func (r *Runner) gotInterrupt() bool {
	select {
	// 当中断事件被触发时发出的信号
	case <-r.interrupt:
		//停止接受后续的任何信号
		signal.Stop(r.interrupt)
		return true
	default:
		return false
	}
}

// run执行每一个已注册的任务
func (r *Runner) run() error {
	for id, task := range r.tasks {
		//检测操作系统的中断信号
		if r.gotInterrupt() {
			return ErrInterrupt
		}

		//执行已注册的任务
		task(id)
	}
	return nil
}

// Start执行所有任务，并监听通道事件
func (r *Runner) Start() error {
	//我们系统接受所有中断信号
	signal.Notify(r.interrupt, os.Interrupt)

	//启动一个goroutine执行任务
	go func() {
		r.complete <- r.run()
	}()

	//这里会发生阻塞，直到有个case能运行
	select {
	//当任务处理完成时发出的信号
	case err := <-r.complete:
		return err
		//当任务处理程序运行超时时发出的信号
	case <-r.timeout:
		return ErrTimeout
	}
}

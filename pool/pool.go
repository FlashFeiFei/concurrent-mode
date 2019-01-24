package pool

import (
	"errors"
	"io"
	"log"
	"sync"
)

//Pool管理一组可以安全地在多个goroutine间
//共享的资源.被管理的资源必须
//实现io.Closer接口
type Pool struct {
	m         sync.Mutex
	resources chan io.Closer
	factory   func() (io.Closer, error)
	closed    bool
}

//ErrPoolClosed请求表示(Acquire)了一个
//已经关闭的池
var ErrPoolClosed = errors.New("Pool has been closed.")

//New创建一个用来管理资源的池
//这个池需要一个可以分配新资源的函数
//并规定池的大小
func New(fn func() (io.Closer, error), size uint) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("size value too small.")
	}
	return &Pool{
		factory:   fn,
		resources: make(chan io.Closer, size),
	}, nil
}

//Acquire从池中获取一个资源
func (p *Pool) Acquire() (io.Closer, error) {
	select {
	//检查是否有空间的资源
	case r, ok := <-p.resources:
		log.Println("Acquire:", "Shared Resource")
		if !ok {
			return nil, ErrPoolClosed
		}
		return r, nil
		//因为没有空闲资源可用，所以提供一个新的资源
	default:
		log.Println("Acquire:", "New Resource")
		return p.factory()
	}
}

//Release将一个使用后的资源放回池里
func (p *Pool) Release(r io.Closer) {
	//保证本操作和close操作的安全
	p.m.Lock()
	defer p.m.Unlock()

	// 如果池已经关闭，销毁这个资源
	if p.closed {
		r.Close()
		return
	}

	select {
	//试图将这个资源放入队列
	case p.resources <- r:
		log.Println("Release:", "In Queue")
		//如果队列已满，则关闭这个资源
	default:
		log.Println("Release:", "Closing")
		r.Close()
	}
}

//Close会让资源池停止工作，并关闭所有现有的资源
func (p *Pool) Close() {
	//保证本操作与Release操作的安全
	p.m.Lock()
	defer p.m.Unlock()

	// 如果pool已经被关闭，什么也不做
	if p.closed {
		return
	}

	//将池关闭
	p.closed = true

	//在清除通道里的资源钱，将通道关闭
	//如果不这样做，会发生死锁
	close(p.resources)

	//将通道里面的资源，一个一个关闭
	for r := range p.resources {
		r.Close()
	}
}

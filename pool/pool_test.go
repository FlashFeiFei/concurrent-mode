package pool

import (
	"io"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	maxGoroutines   = 25 //要使用的goroutine的数量
	pooledResources = 2  //池中的资源的数量
)

//dbConnection模拟要共享的资源
type dbConnection struct {
	ID int32
}

//Close实现了io.Close接口，以便dbConnection
//可以被池管理。Close用来完成任意资源的
//释放管理
func (dbConn *dbConnection) Close() error {
	log.Println("Close:Connection", dbConn.ID)
	return nil
}

//idCounter用来给每个链接分配一个独一无二的id
var idCouter int32

// createConnection是一个工厂函数,
// 当需要一个新连接时，资源池会调用这个函数
func createConnection() (io.Closer, error) {
	id := atomic.AddInt32(&idCouter, 1)
	log.Println("Create:New connection", id)
	return &dbConnection{id}, nil
}

func TestPool(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(maxGoroutines)

	//创建用来管理连接的池
	p, err := New(createConnection, pooledResources)

	if err != nil {
		log.Println(err)
	}

	//使用池里的链接来完成查询
	for query := 0; query < maxGoroutines; query ++ {
		//每个goroutine需要自己复制一份要
		//查询值的副本,不然所有的查询会共享
		//同一个查询变量
		go func(q int) {
			wg.Done()
		}(query)
	}

	//等待goroutine结束
	wg.Wait()

	//关闭池
	log.Println("Shutdown Program.")
	p.Close()
}

func performQueries(query int, p *Pool) {
	//从池里请求一个连接
	conn, err := p.Acquire()
	if err != nil {
		log.Println(err)
		return
	}
	//将该连接释放回池里
	defer p.Release(conn)

	//用等待来模拟查询响应
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Microsecond)
	log.Printf("QID[%d] CID[%d] \n", query, conn.(*dbConnection).ID)
}

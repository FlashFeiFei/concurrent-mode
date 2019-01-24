#无缓冲的通道
无缓冲的通道（unbuffered channel）是指在接收前没有能力保存任何值的通道。*这种类型的通
道要求发送 goroutine 和接收 goroutine 同时准备好，才能完成发送和接收操作*。如果两个 goroutine
没有同时准备好，通道会导致先执行发送或接收操作的 goroutine 阻塞等待。这种对通道进行发送
和接收的交互行为本身就是同步的。其中任意一个操作都无法离开另一个操作单独存在。
**编译代码的时候，无缓冲通道的两个goroutine没有准备好，会出现死锁的报错**


*遍历无缓冲通道*
var court := make(chan int)

for {
    ball , ok := <-court

    if !ok {

    //如果通道被关闭 close(court) 后 ok == false,这个ok表示的是通道被关闭
    //跳出循环
    return
    }
}



#有缓冲的通道
有缓冲的通道（buffered channel）是一种在被接收前能存储一个或者多个值的通道。*这种类
型的通道并不强制要求 goroutine 之间必须同时完成发送和接收*。通道会阻塞发送和接收动作的
条件也会不同。只有在通道中没有要接收的值时，接收动作才会阻塞。只有在通道没有可用缓冲
区容纳被发送的值时，发送动作才会阻塞。这导致有缓冲的通道和无缓冲的通道之间的一个很大
的不同：无缓冲的通道保证进行发送和接收的 goroutine 会在同一时间进行数据交换；有缓冲的
通道没有这种保证


*遍历有缓存通道*

var court := make(chan int)

for {
    ball , ok := <-court

    if !ok {
    //通道关闭后，goroutine 依旧可以从通道接收数据，
    //但是不能再向通道里发送数据。能够从已经关闭的通道接收数据这一点非常重要，因为这允许通
    //道关闭后依旧能取出其中缓冲的全部值，而不会有数据丢失。从一个已经关闭且没有数据的通道
    //里获取数据，总会立刻返回，并返回一个通道类型的零值。如果在获取通道时还加入了可选的标
    //志，就能得到通道的状态信息。
    //这个ok表示通道被关闭并且通道里面没有数据
    return
    }
}
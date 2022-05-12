package main

import (
	"fmt"
	"time"
)

// channel学习

// 工作函数
func worker(c chan int) {
	for n := range c {
		fmt.Printf("Worker received %c\n", n)
	}
}

// 工作函数调用
func CreateWorker(id int) chan<- int {
	c := make(chan int)
	go worker(c)
	return c
}

// 逻辑调用层
func chanDemo() {
	var channels [10]chan<- int
	for i := 0; i < 10; i++ {
		channels[i] = CreateWorker(i) // c 参数是一个管道，需要传输管道数据进去
	}

	for i := 0; i < 10; i++ {
		channels[i] <- 'a' + i
	}

	for i := 0; i < 10; i++ {
		channels[i] <- 'A' + i
	}

	for i := 0; i < 10; i++ {
		channels[i] <- 'V' + i
	}
	time.Sleep(time.Millisecond * 110)
}

// bufferedChannel
func bufferedChannel() chan int {
	c := make(chan int, 3)
	c <- 'a' // 发送数据
	c <- 'b' // 发送数据
	c <- 'c' // 发送数据
	close(c)
	return c
}

// 拿到数据自动关闭
func channelClose() chan int {
	c := make(chan int, 3)
	c <- 'a' // 发送数据
	c <- 'b' // 发送数据
	c <- 'c' // 发送数据
	close(c) // 关闭管道  这样接受数据时就能判断数据是否已经发完了
	return c
}

func main() {
	chanDemo()
	data := channelClose()

	//通过判断数据是否为真来判断是否发完了
	for { // 此时不知道数据是否发完，怎么去知道或者判断数据已经接收完了
		n, ok := <-data // 接收数据  通过ok参数判断是否已经接收完了
		if !ok {
			break
		} else {
			fmt.Printf("%d\n", n)
		}
	}

	// 还有一种方式是通过range的方式获取数据是否已经接收完毕
	bufferData := bufferedChannel()
	for n := range bufferData {
		fmt.Printf("Worker received %c\n", n)
	}
	// 不要通过共享内存来通信， 通过通信来共享内存
}

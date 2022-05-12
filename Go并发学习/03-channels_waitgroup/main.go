package main

import (
	"fmt"
	"sync"
	"time"
)

// channel学习

// 工作函数 done参数来控制数据是否发完
func doWorker(c chan int, w worker) {
	for n := range w.in {
		fmt.Printf("Worker received %c\n", n)
		w.done() // 完成任务
		// wg.Done() // 完成任务

		// go func() {  // 收数据时返回的done 必须要有人接手
		// 	done <- true
		// }()
	}
}

type worker struct {
	in chan int
	// done chan bool
	done func()
}

// 工作函数调用
func CreateWorker(id int, wg *sync.WaitGroup) worker {
	w := worker{
		in: make(chan int),
		done: func() {
			wg.Done()
		},
	}
	go doWorker(w.in, w)
	return w
}

// 逻辑调用层
func chanDemo() {
	// sync.WaitGroup 用来等待所有goroutine执行完毕
	var wg sync.WaitGroup

	var workers [1000000]worker
	for i := 0; i < 1000000; i++ {
		workers[i] = CreateWorker(i, &wg) // c 参数是一个管道，需要传输管道数据进去
	}

	// 添加任务
	wg.Add(2000000)
	for i, worker := range workers {
		worker.in <- 'a' + i
	}
	for i, worker := range workers {
		worker.in <- 'A' + i
	}
	// 等待任务结束
	wg.Wait()
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

// 多任务执行，通过wg.WaitGroup来控制任务是否完成
func main() {
	datanow := time.Now().Unix()
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
	dateEnd := time.Now().Unix()
	fmt.Printf("执行消耗时间: %d 秒", (dateEnd - datanow))
}

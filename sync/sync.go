package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

//互斥锁
//var total struct{
//	sync.Mutex
//	value int
//}

//func worker(wg*sync.WaitGroup) {
//	defer wg.Done()
//	for i := 0; i <= 100; i++ {
//		total.Lock()
//		total.value = total.value + i
//		total.Unlock()
//	}
//}

//atomic包

var total uint64
func worker(wg *sync.WaitGroup) {
	defer wg.Done()
	var i uint64
	for i = 0; i <= 100; i++ {
		atomic.AddUint64(&total, i)
	}
}


//once使用一个数字型标志位，通过原子操作检测标志位状态降低互斥锁的使用
//var once sync.Once
//
//func doSomething(i uint64, total uint64) {
//	once.Do(func() {})
//}

type singleton struct {}
var (
	instance *singleton
	once sync.Once
)

func Instance() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})
	return instance
}

//原子对象数据读写
var config atomic.Value
var x int64
func add(x int64) int64 {
	x++
	return x
}
func doSomething() {
	x := int64(0)
	config.Store(add(x))
	go func() {
		for {
			time.Sleep(time.Second)
			config.Store(add(x))
		}
	}()
	for i:= 0; i < 10; i++ {
		go func() {
			for {
				println(config.Load())
			}
		}()
	}
	time.Sleep(100 * time.Second)
}

//数据一致性模型
var a string
var done bool
func setup() {
	a = "hello world"
	done = true
}

func useChan() {
	done := make(chan int)
	go func() {
		println("hello")
		//done <- 1
		close(done)
	}()
	<-done
}

func helloWithSync() {
	var mu sync.Mutex
	mu.Lock()
	go func() {
		println("helo")
		mu.Unlock()
	}()
	mu.Lock()
}

//生产者消费者模型
func Producer(factor int, out chan<-int) {
	for i := 0; ; i++ {
		out <- i * factor
	}
}

func Consumer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}

func produceAndConsume() {
	ch := make(chan int, 64)
	go Producer(3, ch)
	go Producer(5, ch)
	go Consumer(ch)
	time.Sleep(5 * time.Second)
}

func ProduceAndConsumeHighLevel() {
	ch := make(chan int, 64)
	go Producer(3, ch)
	go Producer(5, ch)
	go Consumer(ch)
	// crtl+c退出
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("quit (%v)\n", <-sig)
}

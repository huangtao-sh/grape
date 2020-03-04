package main

import (
	"fmt"
	"sync"
	"time"
)

func test(wg *sync.WaitGroup) {
	fmt.Println("Hello world")
	time.Sleep(1000)
	fmt.Println("ok2")
	wg.Done()
}
func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go test(&wg)
	go test(&wg)
	wg.Wait()
	fmt.Println("Proce terminated")
}

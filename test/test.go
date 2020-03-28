package main

import (
	"fmt"
	"grape/util"
)

// Print 打印
func Print(data *util.DataCh, w *util.Waiter) {
	defer w.Done()
	for row := range data.Read() {
		fmt.Println(row...)
	}
}

// Send 发送数据
func Send(data *util.DataCh) {
	defer data.Close()
	for i := 0; i < 10; i++ {
		data.Write() <- []interface{}{i, i + 10}
	}
}
func main() {
	waiter := util.NewWaiter()
	defer waiter.Wait()
	data := util.NewDataCh()
	go Print(data, waiter)
	go Send(data)

}

package util

import (
	"fmt"
	"os"
)

// CheckErr 检查是否有错误，并退出操作系统
func CheckErr(err error, exitCode int) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(exitCode)
	}
}

// CheckFatal 检查致命错误
func CheckFatal(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// Dater 数据接口
type Dater interface {
	Next() bool
	Read() []interface{}
}

// Println 打印一行数据
func Println(data Dater) {
	var row []interface{}
	for data.Next() {
		row = data.Read()
		fmt.Println(row...)
	}
}

// Printf 格式打印
func Printf(format string, data Dater) {
	var row []interface{}
	for data.Next() {
		row = data.Read()
		fmt.Printf(format, row...)
	}
}

// Waiter 协和控制
type Waiter struct {
	done chan int
}

// NewWaiter Waiter 构造函数
func NewWaiter() *Waiter {
	done := make(chan int)
	return &Waiter{done}
}

// Done 任务完成
func (w *Waiter) Done() {
	close(w.done)
}

// Wait 等待协程完成
func (w *Waiter) Wait() {
	<-w.done
}

// DataCh 数据通道
type DataCh struct {
	ch chan []interface{}
}

// NewDataCh DataCh 构造函数
func NewDataCh() *DataCh {
	ch := make(chan []interface{})
	return &DataCh{ch}
}

// Read 读取通道
func (d *DataCh) Read() <-chan []interface{} {
	return d.ch
}

// Write 写入通道
func (d *DataCh) Write() chan<- []interface{} {
	return d.ch
}

// Close 关闭通道
func (d *DataCh) Close() {
	close(d.ch)
}

// Data 数据通道，提供
type Data struct {
	DataCh
	Waiter
}

// NewData Data 构造函数
func NewData() *Data {
	return &Data{*NewDataCh(), *NewWaiter()}
}

// PrintlnCh 异步打印
func PrintlnCh(data *Data) {
	defer data.Done()
	for row := range data.Read() {
		fmt.Println(row...)
	}
}

// PrintfCh 异步格式化打印
func PrintfCh(format string, data *Data) {
	defer data.Done()
	for row := range data.Read() {
		fmt.Printf(format, row...)
	}
}

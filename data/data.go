package data

import (
	"fmt"
	"grape/sqlite3"
	"sync"
)

// Data 数据传输
type Data struct {
	wg *sync.WaitGroup    // 进程控制
	ch chan []interface{} // 数据通道
}

// NewData 构造函数
func NewData() *Data {
	ch := make(chan []interface{})
	return &Data{&sync.WaitGroup{}, ch}
}

// Done 执行完成
func (d *Data) Done() {
	d.wg.Done()
}

// Add 增加计数
func (d *Data) Add(count int) {
	d.wg.Add(count)
}

// Wait 阻塞进程，等待协程运行结束
func (d *Data) Wait() {
	d.wg.Wait()
}

// Close 关闭数据通道
func (d *Data) Close() {
	close(d.ch)
}

// ReadCh 读取数据通道
func (d *Data) ReadCh() <-chan []interface{} {
	return d.ch
}

// WriteCh 读取数据通道
func (d *Data) WriteCh() chan<- []interface{} {
	return d.ch
}

// Write 向通道中写入数据
func (d *Data) Write(data ...interface{}) {
	d.ch <- data
}

// Println 打印通道中的数据
func (d *Data) Println() {
	defer d.Done()
	for row := range d.ch {
		fmt.Println(row...)
	}
}

// Printf 打印通道中的数据
func (d *Data) Printf(format string) {
	defer d.Done()
	for row := range d.ch {
		fmt.Printf(format, row...)
	}
}

// Exec 执行 SQL 语句
func (d *Data) Exec(tx *sqlite3.Tx, sql string) (err error) {
	defer d.Done()
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return
	}
	defer stmt.Close()
	for row := range d.ReadCh() {
		_, err = stmt.Exec(row...)
		if err != nil {
			return
		}
	}
	return
}

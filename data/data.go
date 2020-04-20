package data

import "sync"

// Data 数据传输
type Data struct {
	sync.WaitGroup                    // 进程控制
	ch             chan []interface{} // 数据通道
}

// NewData 构造函数
func NewData() *Data {
	ch := make(chan []interface{})
	return &Data{ch: ch}
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

package data

import (
	"fmt"
	"grape/sqlite3"
	"grape/text"
	"grape/util"
	"sync"
)

// Data 数据传输
type Data struct {
	*sync.WaitGroup                    // 进程控制
	ch              chan []interface{} // 数据通道
}

// NewData 构造函数
func NewData() *Data {
	ch := make(chan []interface{}, 32)
	return &Data{&sync.WaitGroup{}, ch}
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
		fmt.Print(util.Sprintf(format, row...))
	}
}

// Exec 执行 SQL 语句 Deprecated: No longer used.
func (d *Data) Exec(tx *sqlite3.Tx, sql string) (err error) {
	return tx.ExecCh(sql, d)
}

// DReader 带通道的数据读取
type DReader struct {
	Reader
	converters []text.ConvertFunc
}

// Reader 数据读取接口
type Reader interface {
	Next() bool
	Read() []string
}

// NewDReader  DReader 构造函数
func NewDReader(r Reader, converters ...text.ConvertFunc) *DReader {
	return &DReader{r, converters}
}

// ReadAll 读取所有数据
func (r *DReader) ReadAll(d *Data) {
	defer d.Close()
	var row []string
	for r.Next() {
		row = r.Read()
		for _, convert := range r.converters {
			row = convert(row)
			if row == nil {
				break
			}
		}
		if row != nil {
			d.Write(text.Slice(row)...)
		}
	}
}

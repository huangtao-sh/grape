package xls

import (
	"grape/text"
	"grape/util"
	"io"

	"github.com/Luxurioust/excelize"
)

// Reader Excel
type Reader struct {
	*excelize.Rows
	converters []text.ConvertFunc
}

// NewRowReader 工作表构造函数
func NewRowReader(xls *excelize.File, sheet string, skip int, converters ...text.ConvertFunc) *Reader {
	rows, err := xls.Rows(sheet)
	util.CheckFatal(err)
	for i := 0; i < skip; i++ {
		rows.Next()
	}
	return &Reader{rows, converters}
}

// NewXlsReader Excel Reader 构造函数
func NewXlsReader(r io.Reader, sheet string, skip int, converters ...text.ConvertFunc) *Reader {
	xls, err := excelize.OpenReader(r)
	util.CheckFatal(err)
	return NewRowReader(xls, sheet, skip, converters...)
}

// ReadAll 读取所有数据
func (r *Reader) ReadAll(d text.Data) {
	defer d.Close()
	var err error
	var row []string
	Row := r.Rows
	for Row.Next() {
		row, err = Row.Columns()
		util.CheckFatal(err)
		for _, convert := range r.converters {
			row = convert(row)
			if row == nil {
				continue
			}
		}
		d.Write(text.Slice(row)...)
	}
}

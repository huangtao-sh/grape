package xls

import (
	"fmt"
	"grape/text"
	"grape"
	"io"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// ConvertDate 转换日期，把 05-16-20 格式的日期转换成 2020-05-16 格式，无法转换则返回原数据
func ConvertDate(d string) string {
	if len(d) == 8 {
		s := strings.Split(d, "-")
		return fmt.Sprintf("20%s-%s-%s", s[2], s[0], s[1])
	}
	return d
}

// Reader Excel
type Reader struct {
	*excelize.Rows
	converters []text.ConvertFunc
}

// NewRowReader 工作表构造函数
func NewRowReader(xls *excelize.File, sheet string, skip int, converters ...text.ConvertFunc) *Reader {
	rows, err := xls.Rows(sheet)
	grape.CheckFatal(err)
	for i := 0; i < skip; i++ {
		rows.Next()
		rows.Columns()
	}
	return &Reader{rows, converters}
}

// NewXlsReader Excel Reader 构造函数
func NewXlsReader(r io.Reader, sheet string, skip int, converters ...text.ConvertFunc) *Reader {
	xls, err := excelize.OpenReader(r)
	grape.CheckFatal(err)
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
		grape.CheckFatal(err)
		for _, convert := range r.converters {
			row = convert(row)
			if row == nil {
				continue
			}
		}
		if row != nil {
			d.Write(text.Slice(row)...)
		}
	}
}

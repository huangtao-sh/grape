package loader

import (
	"grape/util"
	"io"

	"github.com/extrame/xls"
)

// XlsReader Excel 文件读取
type XlsReader struct {
	Sheet   *xls.WorkSheet
	CurLine int
}

// Read 读取当前记录
func (s *XlsReader) Read() (res []string, err error) {
	if s.CurLine > int(s.Sheet.MaxRow) {
		err = io.EOF
	} else {
		row := s.Sheet.Row(s.CurLine)
		for i := row.FirstCol(); i < row.LastCol(); i++ {
			res = append(res, row.Col(i))
		}
		s.CurLine++
	}
	return
}

// NewXlsReader 读取 Excel 文件
func NewXlsReader(filename string, sheet int, skip int) Reader {
	book, err := xls.Open(filename, "utf-8")
	util.CheckFatal(err)
	st := book.GetSheet(sheet)
	return &XlsReader{st, skip}
}

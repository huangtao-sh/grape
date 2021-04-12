package data

import (
	"grape/util"

	"github.com/huangtao-sh/xls"
)

// XlsReader Excel 文件读取
type XlsReader struct {
	Sheet   *xls.WorkSheet
	CurLine uint16
}

// Next 判断是否还有数据
func (s *XlsReader) Next() bool {
	return s.CurLine <= s.Sheet.MaxRow
}

// Read 读取当前记录
func (s *XlsReader) Read() (res []string) {
	row := s.Sheet.Row(int(s.CurLine))
	for i := row.FirstCol(); i < row.LastCol(); i++ {
		res = append(res, row.Col(i))
	}
	s.CurLine++
	return
}

// NewXlsReader 读取 Excel 文件
func NewXlsReader(filename string, sheet int, skip int) Reader {
	book, err := xls.Open(filename, "utf-8")
	util.CheckFatal(err)
	st := book.GetSheet(sheet)
	return &XlsReader{st, uint16(skip)}
}

package loader

import (
	"grape"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/huangtao-sh/xls"
)

// XlsFile 接口
type XlsFile interface {
	Close() // 关闭文件
	Read(sheet int, skip int, converters ...ConvertFunc) Reader
}

// NewXlsFile 新建 xlsFile 根据文件扩展名自动判断
func NewXlsFile(filename string) (file XlsFile, err error) {
	var fp *os.File
	if strings.ToLower(filepath.Ext(filename)) == ".xls" {
		fp, err = os.Open(filename)
		if err != nil {
			return
		}
		book, err := xls.OpenReader(fp, "utf8")
		if err != nil {
			return nil, err
		}
		file = &xlsFile{fp, book}
	} else {
		fp, err = os.Open(filename)
		if err != nil {
			return
		}
		book, err := excelize.OpenReader(fp)
		if err != nil {
			return nil, err
		}
		file = &xlsxFile{fp, book}
	}
	return
}

// xlsFile .xls 文件对 XlsFile 接口实现
type xlsFile struct {
	file *os.File
	book *xls.WorkBook
}

// Close 关闭文件
func (f *xlsFile) Close() {
	f.file.Close()
}

// Read 读取数据
func (f *xlsFile) Read(sheet int, skip int, converters ...ConvertFunc) Reader {
	st := f.book.GetSheet(sheet)
	r := XlsReader{st, skip}
	return NewConverter(&r, converters...)
}

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
		for i := row.FirstCol(); i <= row.LastCol(); i++ {
			res = append(res, row.Col(i))
		}
		s.CurLine++
	}
	return
}

// NewXlsReader 读取 Excel 文件
func NewXlsReader(filename string, sheet int, skip int) Reader {
	book, err := xls.Open(filename, "utf-8")
	grape.CheckFatal(err)
	st := book.GetSheet(sheet)
	return &XlsReader{st, skip}
}

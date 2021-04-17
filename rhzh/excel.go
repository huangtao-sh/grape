package rhzh

import (
	"grape/loader"
	"io"

	"github.com/huangtao-sh/xls"
)

// XlsReader Execl 文件读取
type XlsReader struct {
	fileName   string
	sheet      string
	curLine    int
	//converters []loader.ConvertFunc
	ready      bool
	res        [][]string
}

// NewXlsReader 读取 xls 文件
func NewXlsReader(fileName string, sheet string, skip int) loader.Reader {
	return &XlsReader{fileName: fileName, sheet: sheet, curLine: skip}
}

// Read 读取函数
func (r *XlsReader) Read() (row []string, err error) {
	if !r.ready {
		if xlFile, fn, err := xls.OpenWithCloser(r.fileName, ""); err == nil {
			defer fn.Close()
			r.res, err = xlFile.GetRows(r.sheet)
			if err != nil {
				return nil, err
			}
			r.ready = true
		}
	}
	if r.curLine < len(r.res) {
		row = r.res[r.curLine]
		r.curLine++
	} else {
		err = io.EOF
	}
	return
}

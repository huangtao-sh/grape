package loader

import (
	"fmt"
	"grape/util"
	"io"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// xlsxFile .xlsx .xlsm 文件的 XlsFile 接口实现
type xlsxFile struct {
	file *os.File
	book *excelize.File
}

// Close 关闭文件
func (f *xlsxFile) Close() {
	f.file.Close()
}

type xlsxReader struct {
	*excelize.Rows
}

// Read 读取数据
func (r *xlsxReader) Read() (result []string, err error) {
	if r.Rows.Next() {
		result, err = r.Rows.Columns()
	} else {
		err = io.EOF
	}
	return
}

// Read 读取数据
func (f *xlsxFile) Read(sheet int, skip int, converters ...ConvertFunc) Reader {
	sheetname := f.book.GetSheetName(sheet)
	rows, err := f.book.Rows(sheetname)
	util.CheckFatal(err)
	for i := 0; (i < skip) && rows.Next(); i++ {
		rows.Columns()
	}
	return NewConverter(&xlsxReader{rows}, converters...)
}

// ConvertDate 转换日期，把 05-16-20 格式的日期转换成 2020-05-16 格式，无法转换则返回原数据
func ConvertDate(d string) string {
	if len(d) == 8 {
		s := strings.Split(d, "-")
		return fmt.Sprintf("20%s-%s-%s", s[2], s[0], s[1])
	}
	return d
}

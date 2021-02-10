package xls

import (
	"fmt"
	"grape/path"
	"grape/text"
	"grape/util"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// Cell 把坐标转换成单位元格
func Cell(col interface{}, row int) (cell string) {
	var err error
	switch column := col.(type) {
	case int:
		cell, err = excelize.CoordinatesToCellName(column, row)
		util.CheckFatal(err)
	case string:
		cell, err = excelize.JoinCellName(column, row)
		util.CheckFatal(err)
	}
	return
}

// Write 向 Excel 写入数据 Deprecated
func Write(book *excelize.File, sheet string, axis string, data util.Dater) {
	col, row, err := excelize.CellNameToCoordinates(axis)
	util.CheckFatal(err)
	for ; data.Next(); row++ {
		rowdata := data.Read()
		book.SetSheetRow(sheet, Cell(col, row), &rowdata)
	}
}

// WriteData 向 Excel 写入数据 Deprecated
func WriteData(book *excelize.File, sheet string, axis string, data util.Dater, header string, widthes map[string]float64) {
	if book.GetSheetIndex(sheet) == -1 {
		book.NewSheet(sheet) // 工作表不存在时自动创建
	}
	if widthes != nil {
		SetWidth(book, sheet, widthes) // 设置宽度
	}
	col, row, err := excelize.CellNameToCoordinates(axis) // 读取初始
	util.CheckFatal(err)
	writer, err := book.NewStreamWriter(sheet)
	util.CheckFatal(err)
	if header != "" {
		headers := text.Slice(strings.Split(header, "|"))
		writer.SetRow(Cell(col, row), headers)
		row++
	}
	for ; data.Next(); row++ {
		rowdata := data.Read()
		writer.SetRow(Cell(col, row), rowdata)
	}
	writer.Flush()
}

// SetWidth 设置宽度 Deprecated
func SetWidth(book *excelize.File, sheet string, widthes map[string]float64) {
	for col, width := range widthes {
		aa := strings.Split(col, ":")
		if len(aa) == 1 {
			aa = append(aa, aa[0])
		}
		if len(aa) != 2 {
			panic("单元格格式错")
		}
		book.SetColWidth(sheet, aa[0], aa[1], width)
	}
}

// File Excel 文件
type File struct {
	*excelize.File
}

// NewFile 新建 Excel 文件
func NewFile() *File {
	return &File{excelize.NewFile()}
}

// GetSheet 获取工作表
func (f *File) GetSheet(index interface{}) *WorkSheet {
	var name string
	switch idx := index.(type) {
	case int:
		name = f.GetSheetName(idx)
	case string:
		name = idx
		i := f.GetSheetIndex(name)
		if i == -1 {
			f.NewSheet(name)
		}
	}
	return &WorkSheet{f, name}
}

// SetWidth 设置表格的宽度
func (f *File) SetWidth(sheet string, widthes map[string]float64) {
	for col, width := range widthes {
		aa := strings.Split(col, ":")
		if len(aa) == 1 {
			aa = append(aa, aa[0])
		}
		if len(aa) != 2 {
			panic("单元格格式错")
		}
		f.SetColWidth(sheet, aa[0], aa[1], width)
	}
}

// TableFormat 单位元样式
const TableFormat = `{"table_style":"TableStyleMedium6", "show_first_column":false,"show_last_column":false,"show_row_stripes":true,"show_column_stripes":false}`

// WriteData 设置表格的宽度
func (f *File) WriteData(sheet string, axis string, header string, data util.Dater) {
	var count int
	if f.GetSheetIndex(sheet) == -1 {
		f.NewSheet(sheet) // 工作表不存在时自动创建
	}
	col, row, err := excelize.CellNameToCoordinates(axis) // 读取初始
	util.CheckFatal(err)
	writer, err := f.NewStreamWriter(sheet)
	util.CheckFatal(err)
	if header != "" {
		headers := text.Slice(strings.Split(header, ","))
		count = len(headers)
		writer.SetRow(Cell(col, row), headers)
		row++
	}
	for ; data.Next(); row++ {
		rowdata := data.Read()
		writer.SetRow(Cell(col, row), rowdata)
	}
	end, _ := excelize.CoordinatesToCellName(col+count-1, row-1)
	writer.AddTable(axis, end, TableFormat)
	writer.Flush()
}

// SaveAs 保存文件
func (f *File) SaveAs(p interface{}) {
	var filepath *path.Path
	switch file := p.(type) {
	case string:
		filepath = path.NewPath(file)
	case *path.Path:
		filepath = file
	default:
		panic(fmt.Sprintf("%v不是有效的文件名", f))
	}
	f.File.SaveAs(filepath.String())
}

// WorkSheet 工作表
type WorkSheet struct {
	file *File
	name string
}

// Rename 修改工作表名称
func (s *WorkSheet) Rename(newName string) {
	s.file.SetSheetName(s.name, newName)
	s.name = newName
}

// Write 写入数据
func (s *WorkSheet) Write(axis string, header string, widthes map[string]float64, data util.Dater) {
	if widthes != nil {
		s.file.SetWidth(s.name, widthes)
	}
	s.file.WriteData(s.name, axis, header, data)
}

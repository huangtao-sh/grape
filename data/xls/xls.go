package xls

import (
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

// Write 向 Excel 写入数据
func Write(book *excelize.File, sheet string, axis string, data util.Dater) {
	col, row, err := excelize.CellNameToCoordinates(axis)
	util.CheckFatal(err)
	for ; data.Next(); row++ {
		rowdata := data.Read()
		book.SetSheetRow(sheet, Cell(col, row), &rowdata)
	}
}

// WriteData 向 Excel 写入数据
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
		headers := text.Slice(strings.Split(header, ","))
		writer.SetRow(Cell(col, row), headers)
		row++
	}
	for ; data.Next(); row++ {
		rowdata := data.Read()
		writer.SetRow(Cell(col, row), rowdata)
	}
	writer.Flush()
}

// SetWidth 设置宽度
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

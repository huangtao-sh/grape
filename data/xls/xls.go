package xls

import (
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

// write 向 Excel 写入数据
func Write(book *excelize.File, sheet string, axis string, data util.Dater) {
	col, row, err := excelize.CellNameToCoordinates(axis)
	util.CheckFatal(err)
	for ; data.Next(); row++ {
		rowdata := data.Read()
		book.SetSheetRow(sheet, Cell(col, row), &rowdata)
	}
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

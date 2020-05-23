package main

import (
	"fmt"
	"grape/path"
	"grape/sqlite"
	"grape/util"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func Cell(row, col int) (r string) {
	column, err := excelize.ColumnNumberToName(col)
	util.CheckFatal(err)
	r, err = excelize.JoinCellName(column, row)
	util.CheckFatal(err)
	return
}

func Export(db *sqlite.DB) {
	header := strings.Fields("日期 银联净额（借方） 收付费（借方) 手续费轧差（贷方） 银数轧差（贷方） 记账金额（贷方）")
	sheet := "信用卡报表"
	file := excelize.NewFile()
	file.SetSheetName("Sheet1", sheet)
	file.SetColWidth(sheet, "A", "A", 12)
	file.SetColWidth(sheet, "B", "F", 20)
	style, err := file.NewStyle(`{"number_format":4}`)
	util.CheckFatal(err)
	headerStyle, err := file.NewStyle(`{"alignment":{"horizontal":"center"},"font":{"family":"黑体"}}`)
	util.CheckFatal(err)
	file.SetCellStyle(sheet, "A1", "F1", headerStyle)
	for i, h := range header {
		file.SetCellValue(sheet, Cell(1, i+1), h)
	}
	sql := `select * from (select * from ylb order by rq desc limit 10) order by rq`
	r, err := db.Fetch(sql)
	util.CheckFatal(err)
	var row int
	for row = 2; r.Next(); row++ {
		for i, v := range r.Read() {
			file.SetCellValue(sheet, Cell(row, i+1), v)
		}
	}
	file.SetCellStyle(sheet, "B2", Cell(row-1, 6), style)
	// 获取导出文件名
	date := GetDate(db)
	filename := path.NewPath(fmt.Sprintf("~/信用卡/信用卡报表-%s.xlsx", date))
	file.SaveAs(filename.String())
	fmt.Printf("导出文件 %s 成功！\n", filename)
}

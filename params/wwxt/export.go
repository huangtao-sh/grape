package wwxt

import (
	"fmt"
	"grape/data/xls"
	"grape/sqlite3"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// Export 导出数据主程序
func Export() {
	var date string
	err := sqlite3.QueryRow("select max(date) from wwxt").Scan(&date)
	if err == nil {
		widthes := map[string]float64{
			"A":   10,
			"B":   30,
			"C:E": 15,
		}
		fmt.Println("最新日期：", date)
		book := excelize.NewFile()
		sheet := "新增"
		book.SetSheetName("Sheet1", sheet)
		xls.SetWidth(book, sheet, widthes)
		book.SetSheetRow(sheet, "A1", &[]interface{}{"系统编号", "系统名称", "机构类型", "机构号"})
		rows := sqlite3.Fetch("select id,name,jglx,jgh from wwxt where date=?", date)
		rows.Export(book, sheet, "A2")
		sheet = "历史"
		book.NewSheet(sheet)
		xls.SetWidth(book, sheet, widthes)
		book.SetSheetRow(sheet, "A1", &[]interface{}{"系统编号", "系统名称", "机构类型", "机构号", "维护时间"})
		rows = sqlite3.Fetch("select * from wwxt order by id")
		rows.Export(book, sheet, "A2")
		book.SaveAs(fmt.Sprintf("%s/新增外围系统列表%s.xlsx", Root.String(), date))
		fmt.Println("导出文件成功！")
	} else {
		fmt.Println("无需要导出数据")
	}
}

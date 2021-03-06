package wwxt

import (
	"fmt"
	"grape/data/xls"
	"grape/sqlite3"
)

// Export 导出数据主程序
func Export() {
	var date string
	const (
		Header = "系统编号, 系统名称, 机构类型, 机构号"
	)
	err := sqlite3.QueryRow("select max(date) from wwxt").Scan(&date)
	if err == nil {
		widthes := map[string]float64{
			"A":   10,
			"B":   30,
			"C:E": 15,
		}
		fmt.Println("最新日期：", date)
		book := xls.NewFile()
		sheet := book.GetSheet(0)
		sheet.Rename("新增")
		sheet.Write("A1", Header, widthes, sqlite3.Fetch("select id,name,jglx,jgh from wwxt where date=?", date))
		sheet = book.GetSheet("历史")
		sheet.Write("A1", Header, widthes, sqlite3.Fetch("select * from wwxt order by id"))
		book.SaveAs(fmt.Sprintf("%s/新增外围系统列表%s.xlsx", Root.String(), date))
		fmt.Println("导出文件成功！")
	} else {
		fmt.Println("无需要导出数据")
	}
}

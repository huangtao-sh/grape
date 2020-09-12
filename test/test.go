package main

import (
	"fmt"
	"grape/data/xls"
	"grape/sqlite3"
)

func main() {
	sqlite3.Config("params")
	defer sqlite3.Close()
	//sqlite3.Println("select zh,hm from nbzh limit 10")
	widthes := map[string]float64{"A": 25, "B": 50}
	book := xls.NewFile()
	book.SetWidth("Sheet1", widthes)
	book.WriteData("Sheet1", "A1", "账号,户名", sqlite3.Fetch("select zh,hm from nbzh limit 10"))
	book.SaveAs("~/abc.xlsx")
	fmt.Println("导出文件成功！")
}

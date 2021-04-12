package xyk

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// Export 导出报表
func Export() {
	file := excelize.NewFile()
	s := []interface{}{"hello", 15, true}
	file.NewSheet("Test5")
	file.SetColWidth("Test5", "A", "C", 20)
	st, _ := file.NewStreamWriter("Test5")
	st.SetRow("A1", s)
	st.SetRow("A2", s)
	st.Flush()
	//file.SetActiveSheet(index)
	err := file.SaveAs("c:/users/huangtao/test.xlsx")
	if err != nil {
		println(err.Error())
	}
	println("生成文件成功")
}

package wwxt

import (
	"fmt"
	"grape/sqlite3"
	"grape/util"
	"strings"
)

// Add 新增外围系统名称
func Add(names *string) {
	var No, value int
	err := sqlite3.QueryRow("select max(id) from wwxt").Scan(&No)
	if err != nil {
		No = 0
	}
	No++
	Names := strings.Split(*names, "，")
	fmt.Println("拟新增如下参数：")
	for i, name := range Names {
		err := sqlite3.QueryRow("select id from wwxt where name=?", name).Scan(&value)
		if err != nil {
			fmt.Printf("%3d  %s\n", int(No)+i, name)
		} else {
			fmt.Printf("%s 已存在，编号为 %d \n", name, value)
			return
		}
	}
	fmt.Printf("请确认，Y 或 N ")
	var Confirm string
	fmt.Scanf("%s", &Confirm)
	if (Confirm == "Y") || (Confirm == "y") {
		tx := sqlite3.NewTx()
		defer tx.Rollback()
		stmt, err := tx.Prepare("insert into wwxt(id,name,date) values(?,?,date('now'))")
		util.CheckFatal(err)
		defer stmt.Close()
		for i, name := range Names {
			stmt.Exec(int(No)+i, name)
		}
		tx.Commit()
		fmt.Println("新增参数成功！")
	}
}

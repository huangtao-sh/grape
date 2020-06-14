package km

import (
	"fmt"
	"grape/params"
	"grape/sqlite3"
	"grape/text"
)

func initDb() {
	sqlite3.ExecScript(`
	create table if not exists nbzhmb(
		jglx	text,    	-- 机构类型 
		whrq	text,	 	-- 维护日期	
		km		text,		-- 科目
		bz		text,		-- 币种
		xh		int,		-- 序号
		hmbz	text,		-- 户名标志
		hm		text,		-- 户名
		tzed	real,		-- 透支额度
		cszt	text,		-- 初始状态
		jxbz	text,		-- 计息标志
		primary key(jglx,km,bz,xh)
	)
	`)
}

// Load 导入内部账户模板参数
func Load(path text.File, ver string) {
	initDb()
	loader := params.NewLoader(path, "nbzhmb", ver,
		"insert into nbzhmb values(?,date(?),?,?,?,?,?,?,?,?)")
	err := sqlite3.ExecTx(loader)
	if err == nil {
		fmt.Printf("导入 %s 完成 \n", path.FileInfo().Name())
	} else {
		fmt.Println(err)
	}
	return
}

package km

import (
	"grape/params/load"
	"grape/text"
	"grape/util"
)

var initKmSQL = `
create table if not exists kmzd(
	km  	text primary key,    	-- 科目号
	hzkm	text,		-- 汇总科目
	kmmc	text,		-- 科目名称
	kmjb	text,		-- 科目级别
	jdbz	text,		-- 借贷标志
	kmlx	text,		-- 科目类型
	bz		text		-- 备注
)`
var loadKmSQL = "insert or replace into kmzd values(?,?,?,?,?,?,?)"

// LoadKm 导入科目字典
func LoadKm(file text.File, ver string) {
	r, err := file.Open()
	util.CheckFatal(err)
	defer r.Close()
	reader := text.NewReader(text.Decode(r, false, true), false, text.NewSepSpliter(","),
		text.Include(2, 1, 3, 4, 5, 6, 7))
	loader := load.NewLoader("kmzd", file, ver, reader, initKmSQL, loadKmSQL)
	loader.Load()
}

// Main 科目主函数
func Main() {

}

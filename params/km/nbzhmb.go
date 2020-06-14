package km

import (
	"grape/params/load"
	"grape/text"
	"grape/util"
)

var initSQL = `
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
)`
var loadSQL = "insert into nbzhmb values(?,date(?),?,?,?,?,?,?,?,?)"

// Load 导入内部账户模板参数
func Load(file text.File, ver string) {
	r, err := file.Open()
	util.CheckFatal(err)
	defer r.Close()
	reader := text.NewReader(text.Decode(r, false, true), false, text.NewSepSpliter(","))
	loader := load.NewLoader("nbzhmb", file, ver, reader, initSQL, loadSQL)
	loader.Load()
}

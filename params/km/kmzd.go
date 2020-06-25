package km

import (
	"grape/params/load"
	"grape/text"
	"io"
	"os"
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
func LoadKm(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := text.NewReader(r, false, text.NewSepSpliter(","),
		text.Include(2, 1, 3, 4, 5, 6, 7))
	return load.NewLoader("kmzd", info, ver, reader, initKmSQL, loadKmSQL)
}

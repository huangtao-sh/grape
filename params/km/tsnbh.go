package km

import (
	"grape/data/xls"
	"grape/params/load"
	"grape/text"
	"io"
	"os"
)

var initTsnbh = `
create table if not exists tsnbh(
	code text	primary key,	-- 代码
	name	text,	-- 名称
	zhzl	text,	-- 账户种类
	km		text,	-- 科目
	zdyj	text,	-- 制度依据
	whrq	text	-- 维护日期
)
`
var loadTsnbh = `insert or replace into tsnbh values(?,?,?,?,?,?)`

func conv(s []string) []string {
	s[5] = xls.ConvertDate(s[5])
	return s
}

// LoadTsnbh 导入科目
func LoadTsnbh(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := xls.NewXlsReader(r, "历史参数", 1, text.Include(0, 1, 2, 3, 4, 5), conv)
	return load.NewLoader("tsnbh", info, ver, reader, initTsnbh, loadTsnbh)
}

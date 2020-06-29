package km

import (
	"grape/data/xls"
	"grape/params/load"
	"io"
	"os"
)

var initSxfxm = `
create table if not exists sxfxm(
	gn	text,	-- 功能
	xm	text,	-- 项目
	mx	text,	-- 收费细目
	whrq	text,-- 维护日期
	primary key(xm,mx)
)
`
var loadSxfxm = `insert or replace into sxfxm values(?,?,?,?)`

func conv(s []string) []string {
	if s[0] == "手续费" {
		if len(s) < 4 {
			s = append(s, "")
			s[3] = xls.ConvertDate(s[3])
		}
		return s
	} else {
		return nil
	}
}

// LoadTsnbh 导入科目
func LoadSxfxm(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := xls.NewXlsReader(r, "历史参数", 1, conv)
	return load.NewLoader("sxfxm", info, ver, reader, initSxfxm, loadSxfxm)
}

package km

import (
	"grape/params/load"
	"grape/text"
	"io"
	"os"
)

const (
	initZzzz = `create table if not exists zzzz(
		bh		text  primary key,	-- 自制转账编号
		jglx	text,	-- 机构类型
		czjgh	text,	-- 操作机构号
		bz		text,	-- 币种
		jdbz	text,	-- 借贷标志
		szjglx	text,	-- 所在机构类型
		szjg	text,	-- 所在机构 
		km		text,	-- 科目
		xh		text,	-- 序号
		sfkjg	text,	-- 是否跨机构 
		yxhz	text	-- 允许红字	
	)
`
	loadZzzz = "insert or replace into zzzz values(?,?,?,?,?,?,?,?,?,?,?)"
)

// LoadZzzz 导入定制转账参数
func LoadZzzz(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := text.NewReader(r, false, text.NewSepSpliter(","))
	return load.NewLoader("zzzz", info, ver, reader, initZzzz, loadZzzz)
}

package xxbm

import (
	"grape/params/load"
	"grape/text"
	"io"
	"os"
)

const initSQL = `
create table if not exists xxbm(
	bm		text 	primary key,  -- 编码
	name	text,	-- 名称
	km		text	-- 科目
	)	
`
const loadSQL = `insert into xxbm values(?,?,?)`

// Load 导入文件
func Load(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := text.NewReader(r, false, text.NewSepSpliter(","),
		text.Include(0, 1, 2))
	return load.NewLoader("xxbm", info, ver, reader, initSQL, loadSQL)
}

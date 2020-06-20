package xxbm

import (
	"grape/params/load"
	"grape/text"
	"io"
	"os"
)

var initSQL = `
create table if not exists xxbm(
	bm		text 	primary key,  -- 编码
	name	text,	-- 名称
	km		text	-- 科目
	)	
`
var loadSQL = `insert into xxbm values(?,?,?)`

// Load 导入文件
func Load(info os.FileInfo, r io.Reader, ver string) {
	reader := text.NewReader(r, false, text.NewSepSpliter(","),
		text.Include(0, 1, 2))
	loader := load.NewLoader("xxbm", info, ver, reader, initSQL, loadSQL)
	loader.Load()
}

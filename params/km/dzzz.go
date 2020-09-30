package km

import (
	"grape/params/load"
	"grape/text"
	"io"
	"os"
)

const (
	initDzzz = `create table if not exists dzzz(
		bh	text,	-- 定制转账编号
		xh	text,	-- 定制转账序号
		mc	text,	-- 名称
		czjg	text,	-- 操作机构号
		czjglx	text,	-- 操作机构类型
		czjgfh	text,	-- 操作机构所在分行
		czlwjg	text,	-- 操作机构例外机构
		bz		text,	-- 币种
		jdbz	text,	-- 借贷标志
		zhjg	text,	-- 账户所在机构码
		zhjglx	text,	-- 账户机构类型
		zhjgfh	text,	-- 账户所在分行
		zhlwjg	text,	-- 账户例外机构
		km		text,	-- 科目
		zhxh	int,	-- 序号
		yxkjg	text,	-- 是否允许跨机构
		yxhz	text,	-- 是否允许红字
		primary key(bh,xh)
	)
`
	loadDzzz = "insert or replace into dzzz values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
)

func convDzzz(row []string) []string {
	for i, v := range row {
		if v == "null" {
			row[i] = ""
		}
	}
	return row
}

// LoadDzzz 导入定制转账参数
func LoadDzzz(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := text.NewReader(r, false, text.NewSepSpliter(","), convDzzz)
	return load.NewLoader("dzzz", info, ver, reader, initDzzz, loadDzzz)
}

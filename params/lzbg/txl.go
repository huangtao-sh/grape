package lzbg

import (
	"grape/data/xls"
	"grape/params/load"
	"grape/text"
	"io"
	"os"
)

const initTxl = `
create table if not exists txl(
	br		text,	-- 所在机构 
	dept	text,	-- 部门
	name	text,	-- 姓名
	title	text,	-- 职务
	tel		text,	-- 电话
	fax		text,	-- 传真
	mobile	text,	-- 手机
	email	text
)
`
const loadTxl = `insert into txl values(?,?,?,?,?,?,?,?)`

// LoadTxl 导入通讯录
func LoadTxl(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := xls.NewXlsReader(r, "通讯录", 1, text.Include(0, 1, 2, 3, 4, 5, 6, 7))
	return load.NewLoader("txl", info, ver, reader, initTxl, loadTxl)
}

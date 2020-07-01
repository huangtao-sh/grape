package lzbg

import (
	"grape/data/xls"
	"grape/params/load"
	"grape/text"
	"io"
	"os"
)

var initSQL = `
create table if not exists yyzg(
	gyh		text, 		-- 柜员号
	ygh		text,		-- 员工号
	xm		text,		-- 姓名
	js		text,		-- 角色
	lxdh	text,		-- 联系电话
	mobile	text,		-- 手机
	yx		text,		-- 邮箱
	bz		text,		-- 备注
	jg		text,		-- 机构号
	jgmc	text,		-- 机构名称
	whrq	text,		-- 维护日期
	primary key(gyh,jg)
)
`

var loadSQL = "insert or replace into yyzg values(?,?,?,?,?,?,?,?,?,?,?)"

// LoadYyzg 导入营业主管信息
func LoadYyzg(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := xls.NewXlsReader(r, "Sheet1", 1, text.Include(0, 1, 2, 3, 4, 5, 6, 7, 9, 11, 10))
	return load.NewLoader("yyzg", info, ver, reader, initSQL, loadSQL)
}

var initSXB = `
create table if not exists fhsxb(
	br		text primary key,  -- 分行
	[order]	int					-- 顺序
)
`

var loadSXB = `insert into fhsxb Values(?,?)`

// LoadFhsxb 导入分行顺序表
func LoadFhsxb(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := xls.NewXlsReader(r, "分行顺序表", 1, text.Include(0, 1))
	return load.NewLoader("fhsxb", info, ver, reader, initSXB, loadSXB)
}

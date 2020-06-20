package lzbg

import (
	"grape/data/xls"
	"grape/params/load"
	"grape/path"
	"grape/text"
	"grape/util"
)

var initSQL = `
create table if not exists yyzg(
	gyh		text primary key, -- 柜员号
	ygh		text,		-- 员工号
	xm		text,		-- 姓名
	js		text,		-- 角色
	lxdh	text,		-- 联系电话
	mobile	text,		-- 手机
	yx		text,		-- 邮箱
	bz		text,		-- 备注
	jg		text,		-- 机构号
	jgmc	text,		-- 机构名称
	whrq	text		-- 维护日期
)
`

var loadSQL = "insert or replace into yyzg values(?,?,?,?,?,?,?,?,?,?,?)"

// LoadYyzg 导入营业主管信息
func LoadYyzg(file *path.Path) {
	ver := file.FileInfo().Name()[18:24]
	r, err := file.Open()
	util.CheckFatal(err)
	reader := xls.NewXlsReader(r, "Sheet1", 1, text.Include(0, 1, 2, 3, 4, 5, 6, 7, 9, 11, 10))
	loader := load.NewLoader("yyzg", file.FileInfo(), ver, reader, initSQL, loadSQL)
	loader.Load()
}

var initSXB = `
create table if not exists fhsxb(
	br		text primary key,  -- 分行
	[order]	int					-- 顺序
)
`

var loadSXB = `insert into fhsxb Values(?,?)`

// LoadFhsxb 导入分行顺序表
func LoadFhsxb(file *path.Path) {
	ver := ""
	r, err := file.Open()
	util.CheckFatal(err)
	reader := xls.NewXlsReader(r, "分行顺序表", 1, text.Include(0, 1))
	loader := load.NewLoader("fhsxb", file.FileInfo(), ver, reader, initSXB, loadSXB)
	loader.Load()
	//loader.Test()
}

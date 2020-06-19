package lzbg

import (
	"grape/params/load"
	"grape/path"
	"grape/text"
	"grape/util"
	"strconv"

	"github.com/Luxurioust/excelize"
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

// Reader 营业主管履职报告
type Reader struct {
	file load.File
}

// ReadAll 读取所有数据
func (r *Reader) ReadAll(d text.Data) {
	defer d.Close()
	re, err := r.file.Open()
	util.CheckFatal(err)
	defer re.Close()
	xls, err := excelize.OpenReader(re)
	util.CheckFatal(err)
	rows, err := xls.Rows("Sheet1")
	util.CheckFatal(err)
	for rows.Next() {
		row, _ := rows.Columns()
		_, err := strconv.Atoi(row[0])
		if err == nil {
			d.Write(row[0], row[1], row[2], row[3], row[4], row[5], row[6], row[7], row[9], row[11], row[10])
		}
	}
}

// NewReader 构造函数
func NewReader(file load.File) *Reader {
	return &Reader{file}
}

var loadSQL = "insert or replace into yyzg values(?,?,?,?,?,?,?,?,?,?,?)"

// LoadYyzg 导入营业主管信息
func LoadYyzg(file *path.Path) {
	var path load.File = file
	ver := path.FileInfo().Name()[18:24]
	reader := NewReader(path)
	loader := load.NewLoader("yyzg", path, ver, reader, initSQL, loadSQL)
	loader.Load()
}

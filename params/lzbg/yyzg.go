package lzbg

import (
	"flag"
	"fmt"
	"grape/data/xls"
	"grape/loader"
	"grape/params"
	"grape/params/load"
	"grape/path"
	"grape/sqlite3"
	"grape/text"
	"grape/util"
	"io"
	"os"
	"strings"
)

const initSQL = `
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

const loadSQL = "insert or replace into yyzg values(?,?,?,?,?,?,?,?,?,?,?)"

// LoadYyzg 导入营业主管信息
func LoadYyzg(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := xls.NewXlsReader(r, "Sheet1", 1, text.Include(0, 1, 2, 3, 4, 5, 6, 7, 9, 11, 10))
	return load.NewLoader("yyzg", info, ver, reader, initSQL, loadSQL)
}

func conv(row []string) ([]string, error) {
	s := row[10]
	row[10] = strings.Join([]string{s[:4], s[4:6], s[6:]}, "-")
	return row, nil
}

// LoadZg 导入营业主管信息
func LoadZg(file string) {
	sqlite3.ExecScript(initSQL)
	info := path.NewPath(file).FileInfo()
	reader := loader.NewXlsReader(file, 0, 1)
	reader = loader.NewConverter(reader, loader.Include(2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 11), conv)
	lder := loader.NewLoader(info, "yyzg", loadSQL, reader)
	lder.Ver = util.Extract(`\d*`, info.Name())
	lder.Clear = true
	lder.Check = true
	lder.Load()
	fmt.Println("营业主管表导入成功")
}

const initSXB = `
create table if not exists fhsxb(
	br		text primary key,  -- 分行
	[order]	int					-- 顺序
)
`

const loadSXB = `insert into fhsxb Values(?,?)`

// LoadFhsxb 导入分行顺序表
func LoadFhsxb(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := xls.NewXlsReader(r, "分行顺序表", 1, text.Include(0, 1))
	return load.NewLoader("fhsxb", info, ver, reader, initSXB, loadSXB)
}

// YyzgMain 营业主管查询程序主函数
func YyzgMain() {
	sqlite3.Config("params")
	defer sqlite3.Close()
	var t string
	typ := flag.String("t", "", "主管类型")
	flag.Parse()
	if *typ != "" {
		t = fmt.Sprintf(`and js like "%s%%" `, *typ)
	}
	params.PrintVer("yyzg")
	fmt.Println("工号   姓名       角色            联系电话           手机         机构")
	for _, arg := range flag.Args() {
		if util.FullMatch(`\d{4,9}`, arg) {
			arg = fmt.Sprintf("%s%%", arg)
			sqlite3.Printf("%-6s %-10s %-15s %-15s %11s %-30s\n",
				"select ygh,xm,js,lxdh,mobile,jgmc from yyzg where jg like ? "+t, arg)
		} else {
			arg = fmt.Sprintf("%%%s%%", arg)
			sqlite3.Printf("%-6s %-10s %-15s %-15s %11s %-30s\n",
				"select ygh,xm,js,lxdh,mobile,jgmc from yyzg where(xm like ? or jgmc like ?)"+t, arg, arg)
		}
	}
}

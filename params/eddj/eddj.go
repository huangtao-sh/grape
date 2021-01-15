package eddj

import (
	"fmt"
	"grape/data/xls"
	"grape/params/load"
	"grape/sqlite3"
	"grape/util"
	"io"
	"os"
	"strings"
)

const (
	initSQL = `create table if not exists eddj(
code	text	primary key,	-- 等级代码
name	text,	--	等级名称
ed		text,	--	额度
memo	text	-- 备注
)`
	loadSQL = `insert or replace into eddj values(?,?,?,?)`
)

// conv 转换函数
func conv(row []string) (result []string) {
	if len(row) > 15 && util.FullMatch(`\d{2}`, row[0]) {
		result = make([]string, 4)
		copy(result[:2], row[:2])
		result[2] = strings.Join(row[2:15], "|")
		if row[0] < "51" {
			result[3] = fmt.Sprintf("收：%s，付：%s", row[3], row[5])
		} else {
			result[3] = fmt.Sprintf("转：%s", row[11])
		}
	}
	return
}

//Load 导入函数
func Load(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := xls.NewXlsReader(r, "调整后的额度等级表", 1, conv)
	return load.NewLoader("eddj", info, ver, reader, initSQL, loadSQL)
}

// GetEd 获取额度
func GetEd(code string) (result string) {
	code = fmt.Sprintf("('%s','%s')", string([]byte(code)[:2]), string([]byte(code)[2:]))
	sql := fmt.Sprintf("select group_concat(memo,'，') from eddj where code in %s", code)
	sqlite3.QueryRow(sql).Scan(&result)
	return
}

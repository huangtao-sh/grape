package teller

import (
	"fmt"
	"grape/params/load"
	"grape/text"
	"io"
	"os"
	"strings"
)

const initSQL = `
create table if not exists teller(
    id          text    PRIMARY key,   -- 柜员号
    name        text,   -- 姓名
    telephone   text,   -- 电话
    grade       text,   -- 柜员级别
    [group]     text,   -- 柜组
    branch      text,   -- 机构号
    userid      text,   -- 员工号
    post        text,   -- 岗位
    zxjyz       text,   -- 执行交易组
    zzxe        text,   -- 转账限额
    xjxe        text,   -- 现金限额
    rzlx        text,   -- 认证类型
    zt          text,   -- 状态
    pbjy        text,   -- 屏蔽交易
    gwxz        text,   -- 岗位性质
    qyrq        text,   -- 启用日期
    zzrq        text,   -- 终止日期
    jybz        text,   -- 交易币种
    fqjyz       text,   -- 发起交易组
    zjlx        text,   -- 证件类型
    zjhm        text,   -- 证件号码
    sfyy        text    -- 是否运营人员
)`

const loadSQL = "insert into teller values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

/*
*row[:3], *row[4:8], ','.join(map(str.strip, row[8:-25])),
                row[-25], *row[-23:-20], *row[-10:-3], *row[-2:]
*/

func convert(s []string) (d []string) {
	d = make([]string, 22)
	for i := range s {
		s[i] = strings.TrimSpace(s[i])
	}
	length := len(s)
	copy(d[:3], s[:3])
	copy(d[3:7], s[4:8])
	d[7] = strings.Join(s[8:length-26], ",")
	d[8] = s[length-26]
	copy(d[9:12], s[length-24:length-21])
	copy(d[12:19], s[length-11:length-4])
	copy(d[19:], s[length-3:])
	if len(d[15]) == 8  {
		d[15] = fmt.Sprintf("%s-%s-%s", d[15][:4], d[15][4:6], d[15][6:])
	}
    if len(d[16]) == 8  {
		d[16] = fmt.Sprintf("%s-%s-%s", d[16][:4], d[16][4:6], d[16][6:])
	}
	return
}

// Load 导入文件
func Load(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := text.NewReader(r, false, text.NewSepSpliter(","),
		text.UnQuote, convert)
	return load.NewLoader("teller", info, ver, reader, initSQL, loadSQL)
}

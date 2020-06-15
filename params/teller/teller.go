package teller

import (
	"grape/params/load"
	"grape/text"
	"grape/util"
	"strings"
)

var initSQL = `
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
    zjhm        text    -- 证件号码
)`

var loadSQL = "insert into teller values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

/*
*row[:3], *row[4:8], ','.join(map(str.strip, row[8:-25])),
                row[-25], *row[-23:-20], *row[-10:-3], *row[-2:]
*/

func convert(s []string) (d []string) {
	d = make([]string, 21)
	for i := range s {
		s[i] = strings.TrimSpace(s[i])
	}
	copy(d[:3], s[:3])
	copy(d[3:7], s[4:8])
	d[7] = strings.Join(s[8:60-25], ",")
	d[8] = s[60-25]
	copy(d[9:12], s[60-23:60-20])
	copy(d[12:19], s[60-10:60-3])
	copy(d[19:], s[60-2:])
	return
}

// Load 导入文件
func Load(file text.File, ver string) {
	r, err := file.Open()
	util.CheckFatal(err)
	defer r.Close()
	reader := text.NewReader(text.Decode(r, false, true), false, text.NewSepSpliter(","),
		convert)
	loader := load.NewLoader("teller", file, ver, reader, initSQL, loadSQL)
	loader.Load()
	/*
			d := data.NewData()
			d.Add(1)
			go reader.ReadAll(d)
			go d.Println()
		    d.Wait()
	*/
}

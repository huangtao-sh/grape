package nbzh

import (
	"bufio"
	"grape/params/load"
	"grape/text"
	"grape/util"
	"regexp"
	"strings"
)

var initNbzhSQL = `
create table if not exists nbzh(
	zh text primary key,    --  账号
	jgm text,               --  机构码
	bz text,                --  币种
	hm text,                --  户名
	km text,                --  科目
	yefx text,              --  余额方向 1:借 2:贷 0:两性 记帐以借方为准
	ye real,                --  余额
	qhe real,               --  切换额
	zrye real,              --  昨日余额
	zcll real,              --  正常利率
	fxll real,              --  罚息利率
	fdll real,              --  浮动利率系数
	lxjs real,              --  利息积数
	fxjs real,              --  罚息积数
	qxrq text,              --  起息日期
	khrq text,              --  开户日期
	xhrq text,              --  销户日期
	sbfsr text,             --  上笔发生日期
	mxbs int,               --  明细笔数
	zhzt text,              --  账户状态
	/*
第一位:销户状态
0:未销户
1:已销户
9:被抹帐
第二位:冻结状态
0:未冻结
1:借方冻结
2:贷方冻结
3:双向冻结
第三位:收付现标志
0:不可收付现
1:可收付现
jxbz char 2 N.N 计息标志
第一位:计息方式
0:不计息
1:按月计息
2:按季计息
3:按年计息
第二位:入帐方式
0:计息不入帐
1:计息入帐收息
2:计息入帐付息
	*/
	jxbz text,  -- 计息标志
	sxzh text,  -- 收息账号
	fxzh text,  --  付息账号
	tzed real,  -- 透支额度
	memo text   -- 备注
);
drop view nbzhhz;
create view nbzhhz as 
select b.jglx,a.bz,a.km,cast(substr(a.zh,19,3)as int) as xh,a.hm,sum(abs(a.ye)), 
max(a.sbfsr) from nbzh a 
left join ggjgm b on a.jgm=b.jgm 
where a.zhzt like "0%" 
group by b.jglx,a.km,a.bz,xh;
`

var loadNbzhSQL = `
insert into nbzh values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,date(?),date(?),date(?),date(?),?,?,?,?,?,?,?)
`

func convert(s []string) []string {
	return s[:25]
}

// Load 导入文件
func Load(file text.File, ver string) {
	r, err := file.Open()
	util.CheckFatal(err)
	defer r.Close()
	reader := text.NewReader(text.Decode(r, false, true), false, text.NewSepSpliter(","), convert)
	loader := load.NewLoader("nbzh", file, ver, reader, initNbzhSQL, loadNbzhSQL)
	loader.Load()
	//loader.Test()
}

var initKemuSQL = `
create table if not exists kemu(
	km      text primary key,  -- 科目
	name    text,              -- 名称
	description text           -- 说明
);
`
var loadKemuSQL = `insert or replace into kemu values(?,?,?)`

// KemuReader 科目读取
type KemuReader struct {
	*bufio.Scanner
}

// ReadAll 读取所有数据
func (r *KemuReader) ReadAll(d text.Data) {
	var line string
	KEMU := regexp.MustCompile(`^(\d{4,6})\s*(.*)`)
	BLANK1 := regexp.MustCompile(`第.章。*`)
	BLANK2 := regexp.MustCompile(`本科目为一级科目.*`)
	// AcPattern := regexp.MustCompile(`\d{1,6}`)
	defer d.Close()
	var km, name string
	var sm []string
	for r.Scan() {
		line = strings.TrimSpace(r.Text())
		if line == "" || BLANK1.MatchString(line) || BLANK2.MatchString(line) {
			continue
		} else if KEMU.MatchString(line) {
			if km != "" {
				d.Write(km, name, strings.Join(sm, "\n"))
				sm = nil
			}
			s := KEMU.FindAllStringSubmatch(line, -1)
			km, name = s[0][1], s[0][2]
		} else {
			sm = append(sm, line)
		}
	}
	d.Write(km, name, strings.Join(sm, "\n"))
}

// LoadKemu 导入科目
func LoadKemu(file text.File) {
	Ver := regexp.MustCompile(`\d{6}`)
	ver := Ver.FindString(file.FileInfo().Name())
	r, err := file.Open()
	util.CheckFatal(err)
	defer r.Close()
	r1 := text.Decode(r, false, true)
	reader := &KemuReader{bufio.NewScanner(r1)}
	loader := load.NewLoader("kemu", file, ver, reader, initKemuSQL, loadKemuSQL)
	loader.Load()
	//loader.Test()
}

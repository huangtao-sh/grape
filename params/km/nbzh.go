package km

import (
	"fmt"
	"grape/params/load"
	"grape/sqlite3"
	"grape/text"
	"io"
	"os"
)

const initNbzhSQL = `
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

create table if not exists nbzhhz(
	jglx	text,  -- 机构类型
	bz		text,	-- 币种
	km		text,	-- 科目
	xh		int,	-- 序号
	hm		text,	-- 户名
	ye 		real,	-- 余额
	sbfsr	text	-- 最后发生日
);
create index if not exists nbzhhz_km on nbzhhz(km);
`

const loadNbzhSQL = `
insert into nbzh values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,date(?),date(?),date(?),date(?),?,?,?,?,?,?,?)
`

func convert(s []string) []string {
	return s[:25]
}

// LoadNbzh 导入文件
func LoadNbzh(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := text.NewReader(r, false, text.NewSepSpliter(","), convert)
	return load.NewLoader("nbzh", info, ver, reader, initNbzhSQL, loadNbzhSQL)
}

// CreateNbzhhz 创建内部账户汇总文件
func CreateNbzhhz() {
	var verNbzh, verHz string
	const sql2 = `select ver from LoadFile where name=?`
	sql := `
	insert into nbzhhz  
	select b.jglx,a.bz,a.km,cast(substr(a.zh,19,3)as int) as xh,a.hm,sum(abs(a.ye)), 
	max(a.sbfsr) from nbzh a 
	left join ggjgm b on a.jgm=b.jgm 
	where a.zhzt like "0%" 
	group by b.jglx,a.km,a.bz,xh;`
	tx := sqlite3.NewTx()
	defer tx.Rollback()
	tx.QueryRow(sql2, "nbzh").Scan(&verNbzh)
	tx.QueryRow(sql2, "nbzhhz").Scan(&verHz)
	if verNbzh != verHz {
		tx.Exec(`delete from nbzhhz`)
		tx.Exec(sql)
		tx.Exec(`insert or replace into LoadFile(name,ver) values(?,?)`, "nbzhhz", verNbzh)
		tx.Commit()
		fmt.Println("创建内部账户汇总完成！")
	} else {
		fmt.Println("无需创建内部账户汇总")
	}

}

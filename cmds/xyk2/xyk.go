/*
2020-04-08	对程序进行优化，将收付费纳入手续费
*/
package main

import (
	"flag"
	"fmt"
	"grape/sqlite3"
)

func init() {
	sqlite3.Config("xyk.db")
}

func create_table() {
	sql := `
create table if not exists ylqs(
	rq		text	primary key,  -- 日期
	sffjf	real	default 0,    -- 收付费借方
	sffdf	real	default 0,	  -- 收付费贷方
	jybs	int,	-- 交易笔数
	jyjejf	real,	-- 交易金额借方
	jyjedf	real,	-- 交易金额贷方
	jhfjf	real,	-- 费用借方
	jhfdf	real,	-- 费用贷方
	qsfjf	real,	-- 清算费借方
	qsfdf	real,	-- 清算费贷方
	qsjejf	real,	-- 清算净额借方
	qsjedf	real	-- 清算净额贷方
);
create table if not exists ys(
	rq		text	primary key,  -- 日期
	ysjf	real,	-- 银数借方
	ysdf	real	-- 银数贷方
);
drop view if exists ylb;
create view if not exists ylb as
select a.rq,a.jyjejf,a.sffjf-a.sffdf as sff,a.qsfdf-a.qsfjf+a.jhfdf-a.jhfjf+a.sffdf-a.sffjf as sxf,
b.ysdf-b.ysjf as ys,
b.ysdf-b.ysjf-jhfdf-qsfdf+jhfjf+qsfjf-sffdf+sffjf as jzdf 
from ylqs a 
left join ys b 
on a.rq=b.rq;
	`
	sqlite3.ExecScript(sql)
	// fmt.Println("创建数据库表成功")
}

func version() {
	fmt.Println(`信用卡银联对账程序   Ver 1.0`)
}

func main() {
	var show, auto, init, load, export bool
	var sql string
	flag.BoolVar(&auto, "a", false, "自动执行所有指令")
	flag.BoolVar(&init, "i", false, "初始化数据库")
	flag.BoolVar(&load, "l", false, "导入数据")
	flag.BoolVar(&show, "s", false, "显示数据结果")
	flag.StringVar(&sql, "e", "", "显示数据结果")
	flag.BoolVar(&export, "p", false, "导出报表")
	flag.Parse()
	defer sqlite3.Close()
	version()
	if auto || init {
		create_table()
	}
	if auto || load {
		Load()
	}
	if auto || show {
		Show()
	}
	if auto || export {
		Export()
	}
	if sql != "" {
		sqlite3.Println(sql)
	}
}

func Show() {
	sql := `select * from (select * from ylb order by rq desc limit 10) order by rq`
	fmt.Println("日期       银联净额（借方） 收付费（借方)   手续费轧差（贷方）  银数轧差（贷方）   记账金额（贷方）")
	sqlite3.Printf("%s%19.2f%15.2f%19.2f%19.2f%19.2f\n", sql)
}

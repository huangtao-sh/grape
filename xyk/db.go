package xyk

import (
	"grape/sqlite3"
	"path/filepath"
	"fmt"
)

func Open()(db sqlite3.DB){
	var err error
	path:=filepath.Join(Home,"xyk.db")
	db,err=sqlite3.Open(path)
	if err!=nil{
		panic("打开数据库失败")
	}
	return
}

var CreateSql string =`
create table if not exists qsb(
	rq 		text    primary key,	-- 日期
	jybs	int,					-- 交易笔数
	jyjejf	int,					-- 交易金额借方
	jyjedf	int,					-- 交易金额贷方
	jhfjf	int,					-- 交换费借方
	jhfdf	int,					-- 交换费贷方
	qsfjf	int,					-- 清算费借方
	qsfdf	int,					-- 清算费贷方
	qsjejf	int,					-- 清算净额借方
	qsjedf	int,					-- 清算净额贷方
	sffjf	int default 0,    		-- 收付费借方
	sffdf	int default 0,     		-- 收付费贷方
	ysjf	int,     				-- 银数借方
	ysdf	int      				-- 银数贷方
);	
`

func CreateDB()(err error){
	db:=Open()
	defer db.Close()
	_,err= db.Exec(CreateSql)
	if err!=nil{
		fmt.Println(err)
	}
	return
}

func Query(sql string){
	db:=Open()
	defer db.Close()
	rd, err := db.Fetch(sql)
	if err!=nil{
		fmt.Println(err)
	}
	defer rd.Close()
	for rd.Next() {
		fmt.Println(rd.Read()...)
	}
}
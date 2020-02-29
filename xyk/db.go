package xyk

import (
	"database/sql"
	"grape/sqlite"
	"grape/sqlite3"
	"path/filepath"
)

// Open 打开数据库连接
func Open() (db sqlite3.DB) {
	var err error
	path := filepath.Join(Home, "xyk.db")
	db, err = sqlite3.Open(path)
	CheckErr(err)
	return
}

// CreateSQL 创建数据库语句
var CreateSQL string = `
create table if not exists ylqs(
	rq 		text    primary key,	-- 日期
	jybs	int,					-- 交易笔数
	jyjejf	int,					-- 交易金额借方
	jyjedf	int,					-- 交易金额贷方
	jhfjf	int,					-- 交换费借方
	jhfdf	int,					-- 交换费贷方
	qsfjf	int,					-- 清算费借方
	qsfdf	int,					-- 清算费贷方
	qsjejf	int,					-- 清算净额借方
	qsjedf	int 					-- 清算净额贷方
);
create table if not exists ylsff(
	rq		text	primary key,
	sffjf	int, 		    		-- 收付费借方
	sffdf	int			     		-- 收付费贷方
);
create table if not exists ysqs(
	rq		text 	primary key,
	ysjf	int,     				-- 银数借方
	ysdf	int      				-- 银数贷方
);	

create view if not exists qsb as
select a.*,b.ysjf,b.ysdf,
ifnull(c.sffjf,0),ifnull(c.sffdf,0),
a.qsjejf-a.qsjedf as ylqsje,
a.jhfjf-a.jhfdf+a.qsfjf-a.qsfdf as qsf,
b.ysdf-b.ysjf as ysje,
a.jhfjf-a.jhfdf+a.qsfjf-a.qsfdf-b.ysjf+b.ysdf as jzdf 
from ylqs a 
left join ysqs b on a.rq=b.rq 
left join ylsff c on a.rq=c.rq;

create table if not exists jorj(
	rq		text,  	-- 日期
	kh		text,	-- 卡号
	lsh		text,	-- 流水号
	jfje	int, 	-- 借方金额
	dfje	int 	-- 贷方金额
);


create table if not exists eve(
	rq		text,	-- 日期
	seqno	text,	-- 流水号
	cendt	text,	-- 交易时间
	kh		text,	-- 卡号
	je		int,	-- 交易金额
	jdbz	text,	-- 借贷标志
	qsrq	text,	-- 清算日期
	lsh		text	-- 银数流水号
);


create table if not exists yl(
	rq		text,	-- 日期
	lx		text,	-- 文件类型
	seqno	text,	-- 流水号
	cendt	text,	-- 交易时间
	kh		text,	-- 卡号
	je		int,	-- 交易金额
	bwlx	text,	-- 报文类型
	jylx	text,	-- 交易类型
	oldseq	text,	-- 原交易流水
	olddt	text,	-- 原交易时间
	bz		text,	-- 备注
	zt		text	-- 状态  C-贷记,D-借记,Z-冲正，N-被冲正
);
`

// CreateDB 创建数据库
func CreateDB() (err error) {
	db := Open()
	defer db.Close()
	_, err = db.Exec(CreateSQL)
	CheckErr(err)
	return
}

// Query 执行SQL查询，并输出
func Query(sql string) {
	path := filepath.Join(Home, "xyk.db")
	db, err := sqlite.Open(path)
	CheckErr(err)
	defer db.Close()
	sqlite.ExecQuery(db, sql)
}

// Execer 可以执行语句
type Execer interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Prepare(string) (*sql.Stmt, error)
	Query(string, ...interface{}) (*sql.Rows, error)
}

// Scanner 扫描器
type Scanner interface {
	Scan() bool
	Read() []interface{}
}

// Exec 执行一条 SQL 语句
func Exec(e Execer, sql string, args ...interface{}) {
	_, err := e.Exec(sql, args...)
	CheckErr(err)
}

// ExecMany 执行多条语句
func ExecMany(e Execer, sql string, scanner Scanner) {
	stmt, err := e.Prepare(sql)
	CheckErr(err)
	defer stmt.Close()
	for scanner.Scan() {
		_, err := stmt.Exec(scanner.Read())
		CheckErr(err)
	}
}

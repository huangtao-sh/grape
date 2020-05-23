/*
	绩效考核统计程序

*/
import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

create_sql=`
create table if not exists jxkh(
	jgm 	text,  	-- 机构码
	gyh		text,	-- 柜员号
	jym		text,	-- 交易码
	jybh	text primary key,	-- 交易编号
	jhjbh	text,	-- 叫号机编号
	jyrq	text,	-- 交易日期
	qhsd	text,	-- 取号时点
	jhsd	text,	-- 叫号时点
	ksjysd	text,	-- 开始交易时点
	ywwcsd	text,	-- 业务完成时点
	khdhsj	int,	-- 客户等侯时间
	gyslsj	int,	-- 柜员受理时间
	gyczsj	int,	-- 柜员操作时间
	khywblsj	int	-- 客户业务办理时间
)
`

func init{
	db,error:=Open("sqlite3",":memory:")
	if error==nil{
		panic("连接数据库失败")
	}
}
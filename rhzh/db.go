package rhzh

import (
	"grape/path"
	"grape/sqlite3"
)

const createSQL = `
create table if not exists rhsj(
	zh		text, 	-- 账号
	yhjgdm	text,	-- 银行机构代码
	yhjgmc	text,	-- 银行机构名称
	ckrmc	text,	-- 存款人名称
	zcdqdm	text,	-- 注册地区代码
	hm		text,	-- 账户名称
	zhxz	text,	-- 账户性质
	khxkzh	text,	-- 开户许可证核准号
	khrq	text,	-- 开户日期
	xhrq	text,	-- 销户日期
	zt		text,	-- 账户状态
	bzlx	text,	-- 币种类型
	bzzl	text,	-- 币种性质
	zjxz	text,	-- 资金性质
	zhlb	text	-- 账户类型
);
create index if not exists rhsj_zh on rhsj(zh);
create table if not exists bhsj(
	zh		text primary key, -- 账号
	yshm	text,	-- 户名
	hm		text,	-- 修正户名
	zhlb	text,	-- 账户类别
	khrq	text,	-- 开户日期
	xhrq	text,	-- 销户日期
	zt		text,	-- 状态  
	hdjg	text	-- 核对结果
);
`

func init() {
	path.InitLog()                // 初始化日志
	sqlite3.Config("rhzh")        //设置数据库
	sqlite3.ExecScript(createSQL) // 创建数据库
}

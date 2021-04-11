package lzbg

import (
	"grape/path"
	"grape/sqlite3"
)

// InitSQL 初始化表结构
const InitSQL = `
create table if not exists yyzg(
	gyh		text, 		-- 柜员号
	ygh		text,		-- 员工号
	xm		text,		-- 姓名
	js		text,		-- 角色
	lxdh	text,		-- 联系电话
	mobile	text,		-- 手机
	yx		text,		-- 邮箱
	bz		text,		-- 备注
	jg		text,		-- 机构号
	jgmc	text,		-- 机构名称
	whrq	text,		-- 维护日期
	primary key(gyh,jg)
);

create table if not exists lzbg (
	bt		 	text,	-- 请求标题
	jjcd		text,	-- 紧急程度
	bgr	 		text,	-- 报告人
	ygh			text,	-- 报告人工号
	bm			text,	-- 部门
	jg			text,	-- 机构 
	bgrq		text,	-- 报告日期
	cs			text,	-- 抄送
	fj			text,	-- 附件
	yxqkzc		text,	-- 设备运行情况正常
	yxqkyc		text,	-- 设备运行情况异常
	sbmc		text,	-- 设备名称
	ycnr		text,	-- 异常内容
	spyj		text,	-- 审批意见
	fhjjwt		text,	-- 分行解决
	zhjjwt		text,	-- 总行解决问题
	bfzhjjwt	text,	-- 部分需总行解决问题
	shryj		text,	-- 审核人意见
	fzryj		text,	-- 负责人意见
	bglx		text,	-- 报告人类型
	bgzl		text,	-- 报告种类
	zyx			text,	-- 重要性
	nr			text,	-- 具体内容
	primary key(ygh,bgrq)
);
drop view if exists lz;
create view if not exists lz as 
select distinct bgr,ygh,bglx,jg,substr(bgrq,1,7)as bgq from lzbg
`

// Load 导入数据
func Load() {
	sqlite3.ExecScript(InitSQL)
	file := path.NewPath("~/Downloads").Find("会计履职报告*.xls*")
	LoadLzbg(file)
	file = path.NewPath("~/Downloads").Find("营业主管信息*.xls*")
	LoadYyzg(file)
}

// Main 履职报告主程序
func Main() {
	path.InitLog()
	sqlite3.Config("lzbg")
	defer sqlite3.Close()
	Load()
	report()
}

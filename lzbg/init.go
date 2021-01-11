package lzbg

import (
	"flag"
	"fmt"
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
	bgr	 		text,	-- 报告人
	yzh			text,	-- 报告人工号
	bm			text,	-- 部门
	jg			text,	-- 机构 
	bgrq		text,	-- 报告日期
	cs			text,	-- 抄送
	fj			text,	-- 附件
	yxqk		text,	-- 设备运行情况
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
	primary key(yzh,bgrq)
);
`

// Load 导入数据
func Load() {
	file := path.NewPath("~/Downloads").Find("会计履职报告????-??.xlsx")
	fmt.Println(file)

}

// Main 履职报告主程序
func Main() {
	sqlite3.Config("lzbg")
	defer sqlite3.Close()
	init := flag.Bool("i", false, "初始化数据库")
	load := flag.Bool("l", false, "导入数据")

	flag.Parse()
	if *init {
		sqlite3.ExecScript(InitSQL)
		fmt.Println("初始化数据库完成！")
	}
	if *load {
		Load()
		fmt.Println("导入数据完成！")
	}
}

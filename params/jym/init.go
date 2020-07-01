package jym

import (
	"grape/sqlite3"
	"sync"
)

var createSQL = `
-- 交易码表
CREATE TABLE IF NOT EXISTS jym(
    jym     text primary key,   --交易码
    jymc    text,   -- 交易名称
    jyz     text,   --交易组
    yxj     text,   --优先级
    wdsqjb  text,   --网点授权级别
    zssqjb  text,   --中心授权级别
    wdsq    text,   --网点授权
    zssqjg  text,   --中心授权机构
    zssq    text,   --中心授权
    jnjb    text,   --技能级别
    xzbz    text,   --现转标志
    wb      text,   --外包
    dets    text,   --大额提示
    dzdk    text,   --电子底卡
    sxf     text,   --手续费
    htjc    text,   --后台监测
    szjd    text,   --事中监督
    bssx    text,   --补扫时限
    sc      text,   --审查
    mz      text,   --抹账
    cesq    text,   --超额授权
	fjjyz   text    --辅加交易组
);
-- 事后补扫
create table if not exists shbs (
    jym     primary key
);

-- 磁道校验
create table if not exists cdjy (
    jym     primary key
);
-- 交易组

create table if not exists jyz(
	jyz 	text 	priamry key, -- 交易组
	jyzm	text	-- 交易组名
);

-- 交易菜单
create table if not exists menu(
	jym		text	primary key, -- 交易码
	name	text,	-- 交易名称
	yjcd	text, 	-- 一级菜单
	ejcd	text	-- 二级菜单
);

-- drop view if exists jymb;

create view if not exists jymb as 
select a.jym,jymc,a.jyz,b.jyzm,yxj,wdsqjb,zssqjb,wdsq,zssqjg,zssq,jnjb,xzbz,wb,
dets,dzdk,sxf,htjc,szjd,bssx,sc,mz,cesq,fjjyz,
case when exists(select jym from shbs where jym=a.jym) then "FALSE" else "TRUE" end as shbs,
case when exists(select jym from cdjy where jym=a.jym) then "TRUE" else "FALSE" end as cdjy
from jym a 
left join jyz b on a.jyz=b.jyz;

create view if not exists jycs as 
select a.jymc,a.jym,a.jyz,b.jyzm,yxj,
case wdsqjb when "1" then "1-主办授权" when "2" then "2-主管授权" end as wdsqjb,
case zssqjb when "1" then "1-主办授权" when "2" then "2-主管授权" end as zssqjb,
wdsq,case zssqjg when "0" then "0-总中心" when "1" then "1-分中心" end as zssqjg,
zssq,jnjb,xzbz,
case wb when "1" then "1-不需要" when "2" then "2-需要" end as wb,
case dets when "0" then "0-不需要" when "1" then "1-需要" end as dets,
case dzdk when "0" then "0-不扫描" when "1" then "1-扫描" end as dzdk,
case sxf when "0" then "0-不需要" when "1" then "1-需要" end as sxf,
case htjc when "0" then "0-不需要" when "1" then "1-需要" end as htjc,
case szjd when "0" then "0-不扫描" when "1" then "1-实时扫描" when "2" then "2-补扫" end as sxjd,
bssx,
case sc when "0" then "0-不需要" when "1" then "1-需要" end as sc,
case mz when "0" then "0-不允许" when "1" then "1-允许" end as mz,
cesq,fjjyz,
case when exists(select jym from shbs where jym=a.jym) then "FAlSE" else "TRUE" end as shbs,
case when exists(select jym from cdjy where jym=a.jym) then "TRUE" else "FALSE" end as cdjy,
c.yjcd,c.ejcd
from jym a 
left join jyz b on a.jyz=b.jyz
left join menu c on a.jym=c.jym;
`
var onece = sync.Once{}

// InitJym 初始化交易码相关表结构
func InitJym() {
	onece.Do(
		func() {
			sqlite3.ExecScript(createSQL)
		})
}
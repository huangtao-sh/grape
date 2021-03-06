package ggjgm

import (
	"grape/params/load"
	"grape/text"
	"io"
	"os"
)

const initSQL = `
-- drop table if exists ggjgm;
create table if not exists ggjgm (
    jgm text primary key,	-- 0 机构码
    mc text,		    	-- 1 机构中文名称
    jc text,				-- 3 机构简称  00、01 类型的机构存放分行的简称
    zfhh text,				-- 7 大额支付行号
    jglx text,				--15 机构类型  00-总行清算中心，01-总行营业部，10-分行清算中心，11-分行营业部，12-支行
    kbrq text,  			--16 开办日期
    hzjgm text 				--17 汇总机构码
);

create view if not exists branch as
select jgm,mc,case when substr(jgm,1,2) in ("33","34") then "9"||substr(jgm,2,8)  -- 浙江省机构排最后
when jgm="653000000" then "650000000"                       -- 重庆分行提前
else jgm end as brorder 
from ggjgm where jgm like "%000" and jglx="10" and jgm not in("998930000"); -- 剔除香港分行
`
const loadSQL = "insert or replace into ggjgm values(?,?,?,?,?,date(?),?)"

// Load 导入文件
func Load(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := text.NewReader(r, false, text.NewSepSpliter(","),
		text.Include(0, 1, 3-43, 7-43, 15-43, 16-43, 17-43))
	return load.NewLoader("ggjgm", info, ver, reader, initSQL, loadSQL)
}

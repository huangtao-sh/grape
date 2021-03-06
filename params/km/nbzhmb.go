package km

import (
	"grape/params/load"
	"grape/text"
	"io"
	"os"
)

const initMbSQL = `
create table if not exists nbzhmb(
	jglx	text,    	-- 机构类型 
	whrq	text,	 	-- 维护日期	
	km		text,		-- 科目
	bz		text,		-- 币种
	xh		int,		-- 序号
	hmbz	text,		-- 户名标志
	hm		text,		-- 户名
	tzed	real,		-- 透支额度
	cszt	text,		-- 初始状态
	jxbz	text,		-- 计息标志
	primary key(jglx,km,bz,xh)
)`
const loadMbSQL = "insert into nbzhmb values(?,date(?),?,?,?,?,?,?,?,?)"

// LoadNbzhmb 导入内部账户模板参数
func LoadNbzhmb(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := text.NewReader(r, false, text.NewSepSpliter(","))
	return load.NewLoader("nbzhmb", info, ver, reader, initMbSQL, loadMbSQL)
}

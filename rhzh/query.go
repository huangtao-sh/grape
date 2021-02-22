package rhzh

import (
	"grape/sqlite3"
)

const (
	dbsj = `select a.zh,a.hm,a.zhxz,a.khrq,a.xhrq,a.zt from rhsj a
left join bhsj b on a.zh=b.zh or a.zh="NRA"||b.zh
where b.zh is null
order by a.zh`
	bszq = `select a.zh,a.yshm,b.hm,a.zhlb,b.zhxz,a.khrq,b.khrq,a.xhrq,b.xhrq,a.zt,b.zt 
from bhsj a left join rhsj b on a.zh=b.zh or "NRA"||a.zh=b.zh 
where a.hm=b.hm and a.khrq=b.khrq and a.xhrq=b.xhrq and a.zt=b.zt and a.zhlb=b.zhlb `
)

// Query1 test
func Query1() {
	//sqlite3.Println(bszq)
	sqlite3.Println("select count(zh)from bhsj")
}

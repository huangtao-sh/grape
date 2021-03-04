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
where b.zh is null or a.hm<>b.hm or a.khrq<>b.khrq or a.xhrq<>b.xhrq or a.zt<>b.zt or a.zhlb<>b.zhlb `
)

// Query1 test
func Query1() {
	//sqlite3.Println(bszq)
	sqlite3.Println(bszq)
}

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
where b.zh is null or a.hm<>b.hm or a.khrq<>b.khrq or a.xhrq<>b.xhrq or a.zt<>b.zt or a.zhlb<>b.zhlb`
	cftj = `select a.zh,count(a.zh),sum(case when a.hm=b.hm and a.khrq=b.khrq and a.xhrq=b.xhrq and a.zt=b.zt and a.zhlb=b.zhlb then 1 else 0 end)
from bhsj a left join rhsj b on a.zh=b.zh or "NRA"||a.zh=b.zh 
group by a.zh having count(a.zh)>1`
)

// Query1 test
func Query1() {
	//sqlite3.Println(bszq)
	sqlite3.Println(cftj)
}

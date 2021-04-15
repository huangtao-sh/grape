package rhzh

import (
	"fmt"
	"grape/loader"
	"grape/path"
	"grape/sqlite3"
	"grape/util"
	"log"
	"strings"
)

var (
	acType   map[string]string // 账户类型转换
	acStatus map[string]string // 账户状态转换
	acTypeRh map[string]string // 人行账户类型转换
	loadRhsj string            // 导入人行数据 SQL
	loadBhsj string            // 导入本行数据 SQL
)

func init() {
	acType = map[string]string{
		"结算账户(基本户)":  "基本户",
		"结算账户(一般户)":  "一般户",
		"结算账户(专用户)":  "专用户",
		"电子结算户(一般户)": "一般户",
		"结算账户()":     "结算账户"}

	acTypeRh = map[string]string{
		"基本存款账户":      "基本户",
		"一般存款账户":      "一般户",
		"非预算单位专用存款账户": "专用户",
		"临时机构临时存款账户":  "临时户",
		"预算单位专用存款账户":  "专用户"}

	acStatus = map[string]string{
		"开户":  "正常",
		"销户":  "撤销",
		"不动户": "久悬",
		"待启用": "正常",
		"抹账":  "撤销"}
	loadRhsj = util.Sprintf("insert into rhsj %15V")
	loadBhsj = util.Sprintf("insert or replace into bhsj(zh,khjg,bz,yshm,zhlb,khrq,xhrq,zt,hm) %9V")
}

func convRhsj(row []string) (d []string, err error) {
	d = append(row[:14], acTypeRh[row[6]])
	return
}

// LoadRhsj 导入人行数据
func LoadRhsj() {
	ROOT := path.NewPath("~/Downloads")
	fileName := ROOT.Find("单位银行结算账户开立、变更及撤销情况查询结果输出*.xls")
	fmt.Printf("导入文件:%s\n", fileName)
	file, err := loader.NewXlsFile(fileName)
	if err != nil {
		log.Println("打开文件失败")
	} else {
		defer file.Close()
		reader := file.Read(0, 1, convRhsj)
		lder := loader.NewLoader(path.NewPath(fileName).FileInfo(), "rhsj", loadRhsj, reader)
		err = lder.Load()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("导入文件 %s 成功\n", fileName)
			sqlite3.Printf(
				"导入数据：%d 条\n",
				"select count(*) from rhsj")
		}
	}
}
func convBhsj(row []string) ([]string, error) {
	row = row[:8]
	row[3] = strings.TrimSpace(strings.ToUpper(row[3]))
	row[4] = acType[row[4]]
	row[7] = acStatus[row[7]]
	row = append(row, FullChar(row[3]))
	return row, nil
}

// LoadBhsj 导入本行数据
func LoadBhsj() {
	ROOT := path.NewPath("~/Downloads")
	fileName := ROOT.Find("开户销户登记簿对公账户信息*.xls")
	fmt.Printf("导入文件:%s\n", fileName)
	if file, err := loader.NewXlsFile(fileName); err != nil {
		log.Println("打开文件失败")
	} else {
		defer file.Close()
		reader := file.Read(0, 1, convBhsj)
		lder := loader.NewLoader(path.NewPath(fileName).FileInfo(), "bhsj", loadBhsj, reader)
		lder.Clear = false
		err = lder.Load()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("导入文件 %s 成功\n", fileName)
			sqlite3.Printf(
				"导入数据：%d 条\n",
				"select count(zh) from bhsj")
		}
	}
}

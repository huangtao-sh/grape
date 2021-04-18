package main

import (
	"fmt"
	"grape"
	"grape/loader"
	"grape/sqlite3"
	"log"
	"strings"
)

var (
	acType   map[string]string //本行账户类型转换
	acStatus map[string]string //本行账户状态转换
	acTypeRh map[string]string //人行账户类型转换
	loadRhsj string            //导入人行数据 SQL
	loadBhsj string            //导入本行数据 SQL
)

func init() {
	//本行账户种类转换字典
	acType = map[string]string{
		"结算账户(基本户)":  "基本户",
		"结算账户(一般户)":  "一般户",
		"结算账户(专用户)":  "专用户",
		"电子结算户(一般户)": "一般户",
		"结算账户()":     "结算账户"}
	//人行账户种类转换字典
	acTypeRh = map[string]string{
		"基本存款账户":      "基本户",
		"一般存款账户":      "一般户",
		"非预算单位专用存款账户": "专用户",
		"临时机构临时存款账户":  "临时户",
		"预算单位专用存款账户":  "专用户"}
	//本行账户状态
	acStatus = map[string]string{
		"开户":  "正常",
		"销户":  "撤销",
		"不动户": "久悬",
		"待启用": "正常",
		"抹账":  "撤销"}
	loadRhsj = grape.Sprintf("insert into rhsj %15V")
	loadBhsj = grape.Sprintf("insert or replace into bhsj(zh,khjg,bz,yshm,zhlb,khrq,xhrq,zt,hm) %9V")
}

// convRhsj 人行数据转换程序，新增转换后的账户类型
func convRhsj(row []string) (d []string, err error) {
	d = append(row[:14], acTypeRh[row[6]])
	return
}

// LoadRhsj 导入人行数据
func LoadRhsj() {
	ROOT := grape.NewPath("~/Downloads")
	fileName := ROOT.Find("单位银行结算账户开立、变更及撤销情况查询结果输出*.xls")
	if fileName != "" {
		var err error
		reader := NewXlsReader(fileName, "PAGE1", 1)
		reader = loader.NewConverter(reader, convRhsj)
		info := grape.NewPath(fileName).FileInfo()
		lder := loader.NewLoader(info, "rhsj", loadRhsj, reader)
		err = lder.Load()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("导入文件 %s 成功\n", info.Name())
			sqlite3.Printf(
				"导入数据：%d 条\n",
				"select count(*) from rhsj")
		}
	} else {
		fmt.Println("未在下载目录发现人行数据")
	}
}

//convBhsj 本行数据转换程序
func convBhsj(row []string) ([]string, error) {
	row = row[:8]
	row[3] = strings.TrimSpace(strings.ToUpper(row[3]))
	row[4] = acType[row[4]]
	row[7] = acStatus[row[7]]
	row = append(row, FullChar(row[3]))
	return row, nil
}

//LoadBhsj 导入本行数据
func LoadBhsj() {
	ROOT := grape.NewPath("~/Downloads")
	fileName := ROOT.Find("开户销户登记簿对公账户信息*.xls")
	if fileName != "" {
		var err error
		log.Printf("导入文件：%s", fileName)
		reader := NewXlsReader(fileName, "Sheet1", 1)
		reader = loader.NewConverter(reader, convBhsj)
		info := grape.NewPath(fileName).FileInfo()
		lder := loader.NewLoader(info, "bhsj", loadBhsj, reader)
		lder.Clear = false
		err = lder.Load()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("导入文件 %s 成功\n", info.Name())
			sqlite3.Printf(
				"导入数据：%d 条\n",
				"select count(zh) from bhsj")
		}
	} else {
		fmt.Printf("未在下载目录下发现本行数据")
	}
}

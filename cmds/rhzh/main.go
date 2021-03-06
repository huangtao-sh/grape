package main

import (
	"fmt"
	"grape/rhzh"
)

const (
	description = `人行报送数据比对程序 Version 0.1
用法：
1、将从人行账户管理系统导出的文件 "单位银行结算账户开立、变更及撤销情况查询结果输出" 放在当前用户的下载目录下
2、将从柜面系统 1181-开销户登记簿中导出的数据文件 "开销户登记簿" 放在当前用户的下载目录下
3、运行本程序后将在 下载 目录生成 "账户报送数据比对" 文件`
)

func main() {
	fmt.Println(description)
	rhzh.LoadBhsj()
	rhzh.LoadRhsj()
	rhzh.Query()
}

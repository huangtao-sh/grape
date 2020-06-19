package ggjgm

import (
	"flag"
	"fmt"
	"grape/params"
	"grape/sqlite3"
	"regexp"
	"runtime"
)

// Main 机构码表主函数
func Main() {
	showVer := flag.Bool("V", false, "显示程序版本")
	flag.Parse()
	if *showVer {
		ShowVersion()
	}
	args := flag.Args()
	if len(args) == 0 {
		ShowVersion()
	} else {
		format := "%9s  %-40s  %-10s  %12s  %2s  %10s  %9s\n"
		params.PrintVer("nbzhmb")
		fmt.Println("  机构码   机构名称                                  简称        行号        类型    开立日期  分行")
		for _, arg := range args {
			if matched, _ := regexp.MatchString("316\\d{1,9}", arg); matched {
				arg = fmt.Sprintf("%s%%", arg)
				sqlite3.Printf(format, "select * from ggjgm where zfhh like ?", arg)
			} else if matched, _ := regexp.MatchString("\\d{1,9}", arg); matched {
				arg = fmt.Sprintf("%s%%", arg)
				sqlite3.Printf(format, "select * from ggjgm where jgm like ?", arg)
			} else {
				arg = fmt.Sprintf("%%%s%%", arg)
				sqlite3.Printf(format, "select * from ggjgm where mc like ? ", arg, arg, arg)
			}
		}
	}
}

// ShowVersion 显示程序版本
func ShowVersion() {
	fmt.Printf("Compiled by %s\n", runtime.Version())
	params.PrintVer("nbzhmb")
}

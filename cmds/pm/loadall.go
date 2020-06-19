package main

import (
	"archive/zip"
	"fmt"
	"grape/params/ggjgm"
	"grape/params/km"
	"grape/params/lzbg"
	"grape/params/nbzh"
	"grape/params/teller"
	"grape/path"
	"grape/text"
	"grape/util"
	"sync"
)

// ROOT 参数根目录
var ROOT *path.Path

func init() {
	ROOT = path.NewPath("~/OneDrive/工作/参数备份")
}

// loadFunc 导入函数原型
type ldFunc func(*path.Path)

// Load 导入参数
func Load() {
	var fileList = map[string]ldFunc{}
	fileList[path.NewPath(`~\OneDrive\工作\参数备份\科目说明`).Find("会计科目说明*")] = nbzh.LoadKemu
	fileList[ROOT.Find("运营参数*.zip")] = LoadZip
	fileList[path.NewPath("~/Downloads").Find("营业主管信息*")] = lzbg.LoadYyzg
	wg := &sync.WaitGroup{}
	for file, f := range fileList {
		if file != "" {
			wg.Add(1)
			go func(file string, f ldFunc, wg *sync.WaitGroup) {
				defer wg.Done()
				f(path.NewPath(file))
			}(file, f, wg)
		}
	}
	wg.Wait()
}

// LoadFunc 导入函数类型
type LoadFunc func(text.File, string)

var fileList = map[string]LoadFunc{
	"YUNGUAN_MONTH_STG_ZSRUN_GGJGM.del":    ggjgm.Load,
	"YUNGUAN_MONTH_STG_ZSRUN_GGNBZHMB.del": km.Load,
	"YUNGUAN_MONTH_STG_ZSRUN_GGKMZD.del":   km.LoadKm,
	"users_output.csv":                     teller.Load,
	"YUNGUAN_MONTH_STG_ZSRUN_FHNBHZZ.del":  nbzh.Load,
}

// LoadZip 导入 zip 压缩包
func LoadZip(file *path.Path) {
	ver := file.Base()[12:19]
	fmt.Printf("导入 Zip 参数表，版本号：%s\n", ver)
	f, err := zip.OpenReader(file.String())
	util.CheckFatal(err)
	defer f.Close()
	wwg := &sync.WaitGroup{}
	for _, file := range f.File {
		info := file.FileInfo()
		loadfunc, ok := fileList[info.Name()]
		if ok {
			wwg.Add(1)
			go func(file text.File, ver string, w *sync.WaitGroup) {
				defer w.Done()
				loadfunc(file, ver)
			}(file, ver, wwg)
		}
	}
	wwg.Wait()
	fmt.Printf("文件 %s 已导入\n", file.Base())
}

package main

import (
	"archive/zip"
	"fmt"
	"grape/params/ggjgm"
	"grape/params/km"
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

// Load 导入参数
func Load() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go LoadZip(wg)
	wg.Wait()
}

// LoadFunc 导入函数类型
type LoadFunc func(text.File, string)

var fileList = map[string]LoadFunc{
	"YUNGUAN_MONTH_STG_ZSRUN_GGJGM.del":    ggjgm.Load,
	"YUNGUAN_MONTH_STG_ZSRUN_GGNBZHMB.del": km.Load,
}

// LoadZip 导入 zip 压缩包
func LoadZip(wg *sync.WaitGroup) {
	defer wg.Done()
	zipFile := ROOT.Find("运营参数*.zip")
	ver := path.NewPath(zipFile).Base()[12:19]
	fmt.Printf("导入数据版本号：%s\n", ver)
	f, err := zip.OpenReader(zipFile)
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
}

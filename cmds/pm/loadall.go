package main

import (
	"archive/zip"
	"fmt"
	"grape/params/km"
	"grape/path"
	"grape/util"
	"sync"
)

var ROOT *path.Path

func init() {
	ROOT = path.NewPath("~/OneDrive/工作/参数备份")
}

func Load() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go LoadZip(wg)
	wg.Wait()
}

func LoadZip(wg *sync.WaitGroup) {
	defer wg.Done()
	zipFile := ROOT.Find("运营参数*.zip")
	ver := path.NewPath(zipFile).Base()[12:19]
	fmt.Println(ver)
	f, err := zip.OpenReader(zipFile)
	util.CheckFatal(err)
	defer f.Close()
	wwg := &sync.WaitGroup{}
	for _, file := range f.File {
		info := file.FileInfo()
		switch info.Name() {
		case "YUNGUAN_MONTH_STG_ZSRUN_GGNBZHMB.del":
			wwg.Add(1)
			go km.NewNbzhmbLoader(file, ver).Load(wwg)
		}
	}
	wwg.Wait()
}

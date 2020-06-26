package loadall

import (
	"archive/zip"
	"fmt"
	"grape/gbk"
	"grape/params/ggjgm"
	"grape/params/jym"
	"grape/params/km"
	"grape/params/load"
	"grape/params/lzbg"
	"grape/params/teller"
	"grape/params/xxbm"
	"grape/path"
	"grape/util"
	"io"
	"os"
	"regexp"
	"sync"
)

// ROOT 参数根目录
var ROOT *path.Path

func init() {
	ROOT = path.NewPath("~/OneDrive/工作/参数备份")
}

// LoadFunc 导入函数类型
type LoadFunc func(os.FileInfo, io.Reader, string) *load.Loader

// Load 导入参数
func Load() {
	Ver := regexp.MustCompile(`\d{6,8}`)
	var fileList = map[string]LoadFunc{}
	fileList[ROOT.Find("科目说明/会计科目说明*")] = km.LoadKemu
	fileList[ROOT.Find("营业主管/营业主管信息*")] = lzbg.LoadYyzg
	fileList[ROOT.Find("分行表/分行顺序表*")] = lzbg.LoadFhsxb
	fileList[ROOT.Find("交易菜单/menu*")] = jym.LoadMenu
	//fileList[ROOT.Find("特殊内部账户参数表/特殊内部账户参数*")] = tsnbh.LoadTsnbh
	wg := &sync.WaitGroup{}
	zipfile := ROOT.Find("运营参数*.zip")
	if zipfile != "" {
		wg.Add(1)
		go LoadZip(path.NewPath(zipfile), wg)
	}
	for file, f := range fileList {
		if file != "" {
			wg.Add(1)
			go func(file string, f LoadFunc, wg *sync.WaitGroup) {
				defer wg.Done()
				p := path.NewPath(file)
				info := p.FileInfo()
				ver := Ver.FindString(info.Name())
				r, err := p.Open()
				util.CheckFatal(err)
				defer r.Close()
				f(info, r, ver)
			}(file, f, wg)
		}
	}
	wg.Wait()
}

var fileList = map[string]LoadFunc{
	"YUNGUAN_MONTH_STG_ZSRUN_GGJGM.del":          ggjgm.Load,
	"YUNGUAN_MONTH_STG_ZSRUN_GGNBZHMB.del":       km.LoadNbzhmb,
	"YUNGUAN_MONTH_STG_ZSRUN_GGKMZD.del":         km.LoadKm,
	"users_output.csv":                           teller.Load,
	"YUNGUAN_MONTH_STG_ZSRUN_FHNBHZZ.del":        km.LoadNbzh,
	"YUNGUAN_MONTH_STG_ZSRUN_GGXXBMDZB.del":      xxbm.Load,
	"transactions_output.csv":                    jym.LoadJym,
	"YUNGUAN_MONTH_STG_TELLER_SCANVOUCHER.del":   jym.LoadShbs,
	"YUNGUAN_MONTH_STG_TELLER_TRANSCONTROLS.del": jym.LoadCdjy,
}

// LoadZip 导入 zip 压缩包
func LoadZip(file *path.Path, wwg *sync.WaitGroup) {
	defer wwg.Done()
	ver := file.Base()[12:19]
	fmt.Printf("导入 Zip 参数表，版本号：%s\n", ver)
	f, err := zip.OpenReader(file.String())
	util.CheckFatal(err)
	defer f.Close()
	wg := &sync.WaitGroup{}
	for _, file := range f.File {
		info := file.FileInfo()
		loadfunc, ok := fileList[info.Name()]
		if ok {
			wg.Add(1)
			go func(file *zip.File, ver string, w *sync.WaitGroup) {
				defer w.Done()
				r, err := file.Open()
				util.CheckFatal(err)
				defer r.Close()
				loader := loadfunc(file.FileInfo(), gbk.NewReader(r), ver)
				loader.Load()
			}(file, ver, wg)
		}
	}
	wg.Wait()
	fmt.Printf("文件 %s 处理完毕\n", file.Base())
}

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
	fileList[ROOT.Find("科目说明/会计科目说明*")] = km.LoadKemu         // 会计科目说明
	fileList[ROOT.Find("营业主管/营业主管信息*")] = lzbg.LoadYyzg       // 营业主管
	fileList[ROOT.Find("分行表/分行顺序表*")] = lzbg.LoadFhsxb        // 分行顺序表
	fileList[ROOT.Find("交易菜单/menu*")] = jym.LoadMenu          // 交易菜单
	fileList[ROOT.Find("通讯录/通讯录*")] = lzbg.LoadTxl            // 通讯录
	fileList[ROOT.Find("特殊内部账户参数表/特殊内部账户参数*")] = km.LoadTsnbh // 特殊内部账户参数
	fileList[ROOT.Find("岗位与交易组/岗位及组*")] = jym.LoadJyz         // 交易组
	fileList[ROOT.Find("手续费项目/手续费项目参数*")] = km.LoadSxfxm      // 手续费项目
	zipfile := ROOT.Find("运营参数*.zip")
	if zipfile != "" {
		LoadZip(path.NewPath(zipfile))
	}
	for file, f := range fileList {
		if file != "" {
			func(file string, f LoadFunc) {
				defer util.Recover()
				p := path.NewPath(file)
				info := p.FileInfo()
				ver := Ver.FindString(info.Name())
				r, err := p.Open()
				util.CheckFatal(err)
				defer r.Close()
				loader := f(info, r, ver)
				loader.Load()
			}(file, f)
		}
	}
	km.CreateNbzhhz() // 创建内部账户汇总
}

var fileList = map[string]LoadFunc{
	"YUNGUAN_MONTH_STG_ZSRUN_GGJGM.del":          ggjgm.Load,    // 机构码
	"YUNGUAN_MONTH_STG_ZSRUN_GGNBZHMB.del":       km.LoadNbzhmb, // 内部账户模板
	"YUNGUAN_MONTH_STG_ZSRUN_GGKMZD.del":         km.LoadKm,     // 公共科目字典
	"users_output.csv":                           teller.Load,   // 柜员表
	"YUNGUAN_MONTH_STG_ZSRUN_FHNBHZZ.del":        km.LoadNbzh,   // 内部账户
	"YUNGUAN_MONTH_STG_ZSRUN_GGXXBMDZB.del":      xxbm.Load,     // 公共信息编码
	"transactions_output.csv":                    jym.LoadJym,   // 交易码参数表
	"YUNGUAN_MONTH_STG_TELLER_SCANVOUCHER.del":   jym.LoadShbs,  // 交易码事后补扫
	"YUNGUAN_MONTH_STG_TELLER_TRANSCONTROLS.del": jym.LoadCdjy,  // 交易码磁道校验
	"YUNGUAN_MONTH_STG_TELLER_DZZZCSB.del":       km.LoadDzzz,   // 定制转账参数
	"YUNGUAN_MONTH_STG_TELLER_ZZZZCSB.del":       km.LoadZzzz,   // 定制转账参数

}

// LoadZip 导入 zip 压缩包
func LoadZip(file *path.Path) {
	ver := file.Base()[12:19]
	fmt.Printf("导入 %s ，版本号：%s\n", file, ver)
	f, err := zip.OpenReader(file.String())
	util.CheckFatal(err)
	defer f.Close()
	for _, file := range f.File {
		info := file.FileInfo()
		loadfunc, ok := fileList[info.Name()]
		if ok {
			func(file *zip.File, ver string) {
				defer util.Recover()
				r, err := file.Open()
				util.CheckFatal(err)
				defer r.Close()
				loader := loadfunc(file.FileInfo(), gbk.NewReader(r), ver)
				loader.Load()
			}(file, ver)
		}
	}
	fmt.Printf("文件 %s 处理完毕\n", file.Base())
}

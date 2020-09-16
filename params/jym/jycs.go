package jym

import (
	"fmt"
	"grape/data/xls"
	"grape/date"
	"grape/loader"
	"grape/path"
	"grape/sqlite3"
	"grape/text"
	"grape/util"
	"io"
	"regexp"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// ROOT 交易码参数工作目录
var ROOT *path.Path

func init() {
	ROOT = path.NewPath("~/OneDrive/工作/参数备份/交易码参数")
}

type jycsReader struct {
	r io.Reader
}

func (r *jycsReader) ReadAll(d text.Data) {
	defer d.Close()
	f, err := excelize.OpenReader(r.r)
	util.CheckFatal(err)
	var row []string
	Row, err := f.Rows("新增")
	util.CheckFatal(err)
	Row.Next() // 跳过首行
	Row.Columns()
	for Row.Next() {
		row, err = Row.Columns()
		util.CheckFatal(err)
		row[29] = xls.ConvertDate(row[29])
		row[28] = xls.ConvertDate(row[28])
		dt := text.Slice(row)
		_, err := strconv.Atoi(row[30])
		if err != nil {
			dt[30] = nil
		}
		d.Write(dt...)
	}
}
func convJycs(row []string) []string {
	return row[0:30]
}
func newJycsReader(r io.Reader) loader.Reader {
	return xls.NewXlsReader(r, "新增", 1, convJycs)
}

// LoadJycs 导入交易码参数
func LoadJycs() {
	file := ROOT.Find("交易码参数备份-*")
	Ver := regexp.MustCompile(`\d{6,8}`)
	ver := Ver.FindString(file)
	fmt.Printf("导入文件:%s\n文件版本：%s\n", file, ver)
	loadJycsSQL := util.Sprintf(`insert or replace into jymcs(jymc,jym,jyz,jyzm,yxj,wdsqjb,zxsqjb,wdsq,zxsqjg,   --中心授权机构
			zxsq,jnjb,xzbz,wb,dets,dzdk,sxf,htjc,szjd,bssx,sc,mz,cesq,fjjyz,shbs,cdjy,yjcd,ejcd,bz,cjrq,tcrq) %30V`)
	lder := loader.NewLoader("jymcs", ver, loadJycsSQL, path.NewPath(file), newJycsReader)
	lder.Load()
}

// BackupJycs 导出交易参数
func BackupJycs() {
	file := ROOT.Join(fmt.Sprintf("交易码参数备份-%s.xlsx", date.Today().Format("%Y%M%D")))
	const (
		header = `交易名称,现有系统交易码,交易所属交易组（编码）,交易所属交易组（中文名称）,交易所属优先级（两位数字）,网点交易授权级别:1－主办授权，2－主管授权,中心交易授权级别：1－主办授权，2－主管授权,必须网点授权？TRUE，FALSE,中心授权机构：0-总中心、1-分中心,必须中心授权？TRUE，FALSE,技能级别要求（两位数字）,CashIn：现金收
CashOut：现金付
TransIn：转账收、TransOut：转账付（现转账收与转账付相同）
自助现金收SelfCashIn
自助现金付SelfCashOut
自助转账SelfTransIn/SelfTransOut,是否需要外包：1－不需要，2－需要,是否需要大额提示（大额核查，电话确认）：0－不需要，1－需要,是否需要扫描电子底卡 0-不扫描，1-扫描,是否需要收手续费 0-不需要，1-需要,是否需要后台监测：0－不需要，1－需要,事中监督扫描方式：0-不扫描，1-实时扫描，2-补扫，0-不监督，1-实时监督，2-非实时监督,补扫的限时时间(分钟),是否需要审查（用于调用审查规则的服务）：0－不需要，1－需要,是否允许抹账：0-不允许，1-允许,是否允许超额授权：TRUE－允许，FALSE－不允许,辅助交易组（需与主交易组不一致，以“|”分隔，例：TG001P|TG002P）可为空,是否需要事后补扫TRUE - 需要， FALSE - 不需要,磁道校验信息TRUE - 需要， FALSE - 不需要,一级菜单,二级菜单,备注,创建日期,投产日期,ID`
	)
	fmt.Println("导出备份文件：", file)
	book := xls.NewFile()
	Widthes := map[string]float64{
		"A":     44,
		"B:C":   7,
		"D":     21,
		"E":     7,
		"F:V":   13,
		"W":     22,
		"X:Y":   9,
		"Z":     17,
		"AA":    33,
		"AB":    40,
		"AC:AD": 14,
		"AE":    0,
	}
	book.SetSheetName("Sheet1", "新增")
	book.SetWidth("新增", Widthes)
	book.WriteData("新增", "A1", header, sqlite3.Fetch(`
	select jymc,jym,jyz,jyzm,yxj,wdsqjb,zxsqjb,wdsq,zxsqjg,   --中心授权机构
		zxsq,jnjb,xzbz,wb,dets,dzdk,sxf,htjc,szjd,bssx,sc,mz,cesq,fjjyz,shbs,cdjy,yjcd,ejcd,bz,cjrq,tcrq,rowid from jymcs`))
	book.SaveAs(file)
	fmt.Println("导出备份文件成功！")
}

// Publish 发布最新的投产参数
func Publish() {
	var tcrq string
	today := date.Today().Format("%F")
	fmt.Println(today)
	err := sqlite3.QueryRow("select distinct tcrq from jymcs where tcrq>=? and tcrq <>'' order by tcrq limit 1", today).Scan(&tcrq)
	if err != nil {
		fmt.Println("无待投产的交易")
	} else {
		fmt.Println("投产日期：", tcrq)
	}
}

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
	"os"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const (
	// JycsHeader 交易码数表头
	JycsHeader = `交易名称,现有系统交易码,交易所属交易组（编码）,交易所属交易组（中文名称）,交易所属优先级（两位数字）,网点交易授权级别:1－主办授权，2－主管授权,中心交易授权级别：1－主办授权，2－主管授权,必须网点授权？TRUE，FALSE,中心授权机构：0-总中心、1-分中心,必须中心授权？TRUE，FALSE,技能级别要求（两位数字）,CashIn：现金收
CashOut：现金付
TransIn：转账收、TransOut：转账付（现转账收与转账付相同）
自助现金收SelfCashIn
自助现金付SelfCashOut
自助转账SelfTransIn/SelfTransOut,是否需要外包：1－不需要，2－需要,是否需要大额提示（大额核查，电话确认）：0－不需要，1－需要,是否需要扫描电子底卡 0-不扫描，1-扫描,是否需要收手续费 0-不需要，1-需要,是否需要后台监测：0－不需要，1－需要,事中监督扫描方式：0-不扫描，1-实时扫描，2-补扫，0-不监督，1-实时监督，2-非实时监督,补扫的限时时间(分钟),是否需要审查（用于调用审查规则的服务）：0－不需要，1－需要,是否允许抹账：0-不允许，1-允许,是否允许超额授权：TRUE－允许，FALSE－不允许,辅助交易组（需与主交易组不一致，以“|”分隔，例：TG001P|TG002P）可为空,是否需要事后补扫TRUE - 需要， FALSE - 不需要,磁道校验信息TRUE - 需要， FALSE - 不需要,一级菜单,二级菜单,备注,创建日期,投产日期,ID`
)

var (
	// JycsWidth 交易参数表单位元宽度
	JycsWidth = map[string]float64{
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
	// ROOT 交易码参数工作目录
	ROOT *path.Path
	// Today 当前日期
	Today string
)

func init() {
	ROOT = path.NewPath("~/OneDrive/工作/参数备份/交易码参数")
	Today = fmt.Sprint(date.Today())
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

func newJycsReader(r io.Reader) loader.Reader {
	return xls.NewXlsReader(r, "交易参数备份", 1)
}

// LoadJycs 导入交易码参数
func LoadJycs() {
	file := ROOT.Find("交易码参数备份-????-??-??.*")
	ver := util.Extract(`\d{6,8}`, file)
	fmt.Printf("导入文件:%s\n文件版本：%s\n", file, ver)
	loadJycsSQL := util.Sprintf(`insert or replace into jymcs(jymc,jym,jyz,jyzm,yxj,wdsqjb,zxsqjb,wdsq,zxsqjg,   --中心授权机构
			zxsq,jnjb,xzbz,wb,dets,dzdk,sxf,htjc,szjd,bssx,sc,mz,cesq,fjjyz,shbs,cdjy,yjcd,ejcd,bz,cjrq,tcrq) %30V`)
	lder := loader.NewLoader("jymcs", ver, loadJycsSQL, path.NewPath(file), newJycsReader)
	lder.Load()
}

// BackupJycs 导出交易参数
func BackupJycs() {
	const sheet = "交易参数备份"
	file := ROOT.Join(fmt.Sprintf("交易码参数备份-%s.xlsx", Today))
	fmt.Println("导出备份文件：", file)
	book := xls.NewFile()
	book.SetSheetName("Sheet1", sheet)
	book.SetWidth(sheet, JycsWidth)
	book.WriteData(sheet, "A1", JycsHeader, sqlite3.Fetch(`
select jymc,jym,jyz,jyzm,yxj,wdsqjb,zxsqjb,wdsq,zxsqjg,   --中心授权机构
zxsq,jnjb,xzbz,wb,dets,dzdk,sxf,htjc,szjd,bssx,sc,mz,cesq,fjjyz,shbs,cdjy,yjcd,ejcd,bz,cjrq,tcrq,rowid from jymcs`))
	book.SaveAs(file)
	fmt.Println("导出备份文件成功！")
}

// Publish 发布最新的投产参数
func Publish() {
	const (
		header = `交易名称,现有系统交易码,交易所属交易组（编码）,交易所属交易组（中文名称）,交易所属优先级（两位数字）,网点交易授权级别:1－主办授权，2－主管授权,中心交易授权级别：1－主办授权，2－主管授权,必须网点授权？TRUE，FALSE,中心授权机构：0-总中心、1-分中心,必须中心授权？TRUE，FALSE,技能级别要求（两位数字）,CashIn：现金收
CashOut：现金付
TransIn：转账收、TransOut：转账付（现转账收与转账付相同）
自助现金收SelfCashIn
自助现金付SelfCashOut
自助转账SelfTransIn/SelfTransOut,是否需要外包：1－不需要，2－需要,是否需要大额提示（大额核查，电话确认）：0－不需要，1－需要,是否需要扫描电子底卡 0-不扫描，1-扫描,是否需要收手续费 0-不需要，1-需要,是否需要后台监测：0－不需要，1－需要,事中监督扫描方式：0-不扫描，1-实时扫描，2-补扫，0-不监督，1-实时监督，2-非实时监督,补扫的限时时间(分钟),是否需要审查（用于调用审查规则的服务）：0－不需要，1－需要,是否允许抹账：0-不允许，1-允许,是否允许超额授权：TRUE－允许，FALSE－不允许,辅助交易组（需与主交易组不一致，以“|”分隔，例：TG001P|TG002P）可为空,是否需要事后补扫TRUE - 需要， FALSE - 不需要,磁道校验信息TRUE - 需要， FALSE - 不需要`
	)
	var tcrq string
	today := date.Today().Format("%F")
	fmt.Println(today)
	err := sqlite3.QueryRow("select distinct tcrq from jymcs where tcrq>=? and tcrq <>'' order by tcrq limit 1", today).Scan(&tcrq)
	if err != nil {
		fmt.Println("无待投产的交易")
	} else {
		fmt.Println("投产日期：", tcrq)
		book := xls.NewFile()
		sheet := book.GetSheet("Sheet1")
		sheet.Rename("新增交易码")
		sheet.Write("A1", header, JycsWidth, sqlite3.Fetch(`select jymc,jym,jyz,jyzm,yxj,wdsqjb,zxsqjb,wdsq,zxsqjg,zxsq,jnjb,xzbz,wb,dets,dzdk,sxf,htjc,szjd,bssx,sc,mz,cesq,fjjyz,shbs,cdjy from jymcs where tcrq=?`, tcrq))
		sheet = book.GetSheet("事后监督参数")
		sheet.Write("A1",
			"交易名称,交易码,总行审查,分行审查,流水勾对",
			map[string]float64{
				"A": 44,
				"B": 8.47,
			},
			sqlite3.Fetch("select jymc,jym,'0','0','0' from jymcs where tcrq=?", tcrq),
		)
		sheet = book.GetSheet("绩效-系统交易码")
		sheet.Write("A1",
			"交易名称,交易码,交易类型（内部/外部/其他）",
			map[string]float64{
				"A": 44,
				"B": 8.47,
				"C": 27,
			},
			sqlite3.Fetch("select jymc,jym,'外部' from jymcs where tcrq=?", tcrq),
		)
		sheet = book.GetSheet("绩效-交易折算系数")
		sheet.Write("A1",
			"交易名称,交易码,折算系数",
			map[string]float64{
				"A": 44,
				"B": 8.47,
				"C": 10,
			},
			sqlite3.Fetch("select jymc,jym,1 from jymcs where tcrq=?", tcrq),
		)
		book.SaveAs(ROOT.Join(fmt.Sprintf("交易码参数%s.xlsx", tcrq)))
		fmt.Printf("导出交易码参数文件：%s\n", fmt.Sprintf("交易码参数%s.xlsx", tcrq))
	}
}

// UpdateJycs 更新交易码参数
func UpdateJycs() {
	const (
		sheetName = "交易码参数"
		exportSQL = `select *,rowid from jymcs where(tcrq="" or tcrq is null or tcrq>=?)and jym  not in (select jym from jym)`
		checkSQL  = "select count(name) from LoadFile where name=? and path=? and mtime>=datetime(?)"
		doneSQL   = "insert or replace into LoadFile values(?,?,datetime(?),?)"
	)
	var (
		info  os.FileInfo
		count int
	)
	fmt.Println("更新交易码参数")
	p := ROOT.Join("交易码参数.xlsx")
	tx := sqlite3.NewTx()
	defer tx.Rollback()
	if p.IsExist() {
		info = p.FileInfo()
		err := tx.QueryRow(checkSQL, "jymcs", info.Name(), info.ModTime()).Scan(&count)
		util.CheckFatal(err)
		if count > 0 {
			fmt.Printf("文件 %s 已导入，忽略\n", info.Name())
			return
		}
		fmt.Printf("导入文件：%s\n", p)
		upJycs(p.String(), tx)
	} else {
		fmt.Printf("文件 %s 不存在，直接导出\n", p)
	}
	book := xls.NewFile()
	book.SetSheetName("Sheet1", sheetName)
	book.SetWidth(sheetName, JycsWidth)
	book.WriteData(sheetName, "A1", JycsHeader, tx.Fetch(exportSQL, Today))
	book.SaveAs(p.String())
	fmt.Println("导出交易码参数成功！")
	info = p.FileInfo()
	tx.Exec(doneSQL, "jymcs", info.Name(), info.ModTime(), "1.0")
	tx.Commit()
}

// upJycs 更新交易码参数
func upJycs(file string, tx *sqlite3.Tx) {
	var (
		row       []string
		deleteSQL string = "delete from jymcs where rowid=?"
		insertSQL string = util.Sprintf(`insert or replace into jymcs(jymc,jym,jyz,jyzm,yxj,wdsqjb,zxsqjb,wdsq,zxsqjg,zxsq,jnjb,xzbz,wb,dets,dzdk,sxf,htjc,szjd,bssx,sc,mz,cesq,fjjyz,shbs,cdjy,yjcd,ejcd,bz,cjrq,tcrq,rowid) %31V`)
	)
	book, err := excelize.OpenFile(file)
	util.CheckFatal(err)
	rows, err := book.Rows("交易码参数")
	util.CheckFatal(err)
	for i := 0; i < 1; i++ {
		rows.Next()
		rows.Columns()
	}
	for rows.Next() {
		row, err = rows.Columns()
		util.CheckFatal(err)
		if len(row) > 28 && row[0] != "" {
			for len(row) < 31 {
				row = append(row, "")
			}
			s := text.Slice(row)
			s[28] = xls.ConvertDate(row[28])
			s[29] = xls.ConvertDate(row[29])
			if row[30] == "" {
				s[30] = nil
			}
			if len(row) == 31 && row[29] == "删除" && row[30] != "" {
				_, err = tx.Exec(deleteSQL, row[30])
				util.CheckFatal(err)
				fmt.Printf("删除交易 %s:%s-%s\n", row[30], row[1], row[0])
			} else {
				_, err = tx.Exec(insertSQL, s...)
				util.CheckFatal(err)
			}
		}
	}
}

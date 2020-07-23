package jym

import (
	"grape/data/xls"
	"grape/loader"
	"grape/text"
	"grape/util"
	"io"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

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
	return &jycsReader{r}
}

// LoadJycs 导入交易码参数
func LoadJycs(file loader.File, ver string) *loader.Loader {
	loadJycsSQL := util.Sprintf(`insert or replace into jymcs(jymc,jym,jyz,jyzm,yxj,wdsqjb,zxsqjb,wdsq,zxsqjg,   --中心授权机构
		zxsq,jnjb,xzbz,wb,dets,dzdk,sxf,htjc,szjd,bssx,sc,mz,cesq,fjjyz,shbs,cdjy,yjcd,ejcd,bz,cjrq,tcrq,rowid) %31V`)
	return loader.NewLoader("jymcs", ver, loadJycsSQL, file, newJycsReader)
}


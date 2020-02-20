package xyk

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
)

// LoadRd1002 读取 RD1002 文件
func (r *Reader) LoadRd1002(tx Execer) {
	data, _ := ReadAll(r.rd1002, true)
	buf := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 10 && fields[0] == "收付费" {
			Exec(tx, "insert or replace into ylsff values(?,?,?)", r.date, Atoi(fields[2]), Atoi(fields[3]))
		} else if len(fields) == 11 && fields[0] == "总" {
			Exec(tx, "insert or replace into ylqs values(?,?,?,?,?,?,?,?,?,?)", r.date,
				Atoi(fields[2]), Atoi(fields[3]), Atoi(fields[4]), Atoi(fields[5]), Atoi(fields[6]), Atoi(fields[7]), Atoi(fields[8]), Atoi(fields[9]), Atoi(fields[10]))
		}
	}
}

// LoadInds 读取银联明细数据
func (r *Reader) LoadInds(tx Execer) {
	for _, path := range r.indfiles {
		lx := filepath.Base(path.Name)[11:]
		switch lx {
		case "ACOMA":
			r.LoadAcoma(path, tx)
		case "ICOMN":
			r.LoadICOMN(path, tx)
		default:
			fmt.Println("Skiped", path.Name)
		}
	}
}

// COMA coma 文件分割
var COMA = []int{0, 11, 23, 30, 41, 61, 74, 88, 100, 105, 112, 117, 126, 142,
	155, 158,
	165, 177, 184, 187, 191, 204, 217, 230, 232, 236, 238, 240, 251, 263, 265,
	268, 272, 285, 299, 311, 327, 332, 335, 376, 381, 390, 413, 422, 434, 457,
	479, 488, 500, 523, 545, 567, 589, 611, 621, 631, 641, 650, 673, 714}

// COMAR coma
var COMAR = []int{2, 3, 57, 5, 8, 9, 17, 27}

/*
'lsh': 2,  # 系统跟踪号
'jysj': 3,  # 交易传输时间
'kh': 4,  # 主账号
'jyje': (5, jine),  # 交易金额
'cdje': (6, jine),  # 部分代收时的承兑金额
'sxf': (7, jine),  # 持卡人手续费
'bwlx': 8,  # 报文类型
'jylx': 9,  # 交易类型
'fwdtjm': 14,  # 服务点条件码
'sqydm': 15,  # 授权应答码
'ylsh': 17,  # 原交易系统跟踪号
'ysjhf': (20, jine),  # 应收交换费
'yfjhf': (21, jine),  # 应付交换费
'qsf': (22, jine),  # 转接清算费
'yjysj': 27,  # 原交易时间
'zrkh': -5,      # 转入卡号

*/

// Split 拆分数据
func Split(bytes []byte, list []int) (result []string) {
	start := list[0]
	for _, end := range list[1:] {
		result = append(result, strings.TrimSpace(string(bytes[start:end])))
		start = end
	}
	return
}

// LoadAcoma 读取 ACOMA 文件
func (r *Reader) LoadAcoma(path *zip.File, tx Execer) {
	data, _ := ReadAll(path, false)
	buf := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buf)
	Exec(tx, "delete from yl where lx=? and rq=?", "ACOMA", r.date)
	stmt, err := tx.Prepare("insert into yl(rq,lx,seqno,cendt,kh,je,bwlx,jylx,oldseq,olddt,jdbz)values(?,?,?,?,?,?,?,?,?,?,'C')")
	if err != nil {
		panic("准备 SQL 语句失败")
	}
	for scanner.Scan() {
		s := Split(scanner.Bytes(), COMA)
		data := []interface{}{r.date, "ACOMA"}
		for _, i := range COMAR {
			data = append(data, s[i])
		}
		//fmt.Println(data...)
		stmt.Exec(data...)
	}
}

// LoadICOMN 读取 ACOMA 文件
func (r *Reader) LoadICOMN(path *zip.File, tx Execer) {
	COMNR := []int{2, 3, 4, 5, 8, 9, 17, 27}
	data, _ := ReadAll(path, false)
	buf := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buf)
	Exec(tx, "delete from yl where lx=? and rq=?", "ICOMN", r.date)
	//stmt, err := tx.Prepare("insert into yl(rq,lx,seqno,cendt,kh,je,bwlx,jylx,oldseq,olddt,jdbz)values(?,?,?,?,?,?,?,?,?,?,'C')")
	//if err != nil {
	//	panic("准备 SQL 语句失败")
	//}
	for scanner.Scan() {
		s := Split(scanner.Bytes(), COMA[:30])
		data := []interface{}{r.date, "ICOMN"}
		for _, i := range COMNR {
			data = append(data, s[i])
		}
		if data[6] == "0420" {
			fmt.Println(data...)
		}
		//stmt.Exec(data...)
	}
}

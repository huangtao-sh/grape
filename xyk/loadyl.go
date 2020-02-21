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
			r.LoadACOMA(path, tx)
		case "ICOMN":
			r.LoadICOMN(path, tx)
		case "IERR":
			r.LoadIERR(path, tx)
		case "ITFL":
			r.LoadITFL(path, tx)
		case "AFCP":
			r.LoadAFCP(path, tx)
		case "IFCP":
			r.LoadAFCP(path, tx)
		default:
			fmt.Println("Skiped", path.Name)
		}
	}
}

// ChongZheng 处理冲正的报文
func (r *Reader) ChongZheng(tx Execer) {
	rq := r.date
	var seqno, cendt string
	rows, err := tx.Query("select oldseq,olddt from yl where rq=? and bwlx='0420' ", rq)
	if err != nil {
		panic("执行查询失败")
	}
	for rows.Next() {
		rows.Scan(&seqno, &cendt)
		Exec(tx, "update yl set zt='N' where rq=? and seqno=? and cendt=?  ", rq, seqno, cendt)
	}

}

// LoadAFCP 导入 AFCP 文件
func (r *Reader) LoadAFCP(path *zip.File, tx Execer) {
	offsets := []int{0, 11, 23, 30, 41, 61, 74, 88, 100, 105, 112, 117, 126, 142, 155,
		158, 165, 177, 184, 187, 191, 204, 217, 230, 232, 236, 238, 240,
		251, 263, 265, 268, 272, 285, 371, 373, 413}
	columns := []int{2, 3, 4, 5, 8, 9, 17, 27}
	data, _ := ReadAll(path, false)
	buf := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buf)

	Exec(tx, "delete from yl where lx=? and rq=?", "AFCP", r.date)
	stmt, err := tx.Prepare("insert into yl(rq,lx,seqno,cendt,kh,je,bwlx,jylx,oldseq,olddt,zt)values(?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic("准备 SQL 语句失败")
	}
	for scanner.Scan() {
		s := SplitData(scanner.Bytes(), offsets, columns)
		data := []interface{}{r.date, "AFCP"}
		for _, i := range s {
			data = append(data, i)
		}
		zt := "D"
		data = append(data, zt)
		stmt.Exec(data...)
		fmt.Println(data...)
	}
	fmt.Println("导入 AFCP 文件完成 ")
}

// COMA coma 文件分割
var COMA = []int{0, 11, 23, 30, 41, 61, 74, 88, 100, 105, 112, 117, 126, 142, 155, 158,
	165, 177, 184, 187, 191, 204, 217, 230, 232, 236, 238, 240, 251, 263, 265,
	268, 272, 285, 299, 311, 327, 332, 335, 376, 381, 390, 413, 422, 434, 457,
	479, 488, 500, 523, 545, 567, 589, 611, 621, 631, 641, 650, 673, 714}

// COMAR coma
var COMAR = []int{2, 3, 57, 5, 8, 9, 17, 27}

// Split 拆分数据
func Split(bytes []byte, list []int) (result []string) {
	start := list[0]
	for _, end := range list[1:] {
		result = append(result, strings.TrimSpace(string(bytes[start:end])))
		start = end
	}
	return
}

// LoadACOMA 读取 ACOMA 文件
func (r *Reader) LoadACOMA(path *zip.File, tx Execer) {
	data, _ := ReadAll(path, false)
	buf := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buf)
	Exec(tx, "delete from yl where lx=? and rq=?", "ACOMA", r.date)
	stmt, err := tx.Prepare("insert into yl(rq,lx,seqno,cendt,kh,je,bwlx,jylx,oldseq,olddt,zt)values(?,?,?,?,?,?,?,?,?,?,'C')")
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

// LoadICOMN 读取 ICOMN 文件
func (r *Reader) LoadICOMN(path *zip.File, tx Execer) {
	COMNR := []int{2, 3, 4, 5, 8, 9, 17, 27}
	data, _ := ReadAll(path, false)
	buf := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buf)
	Exec(tx, "delete from yl where lx=? and rq=?", "ICOMN", r.date)
	stmt, err := tx.Prepare("insert into yl(rq,lx,seqno,cendt,kh,je,bwlx,jylx,oldseq,olddt,zt)values(?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic("准备 SQL 语句失败")
	}
	for scanner.Scan() {
		s := Split(scanner.Bytes(), COMA[:30])
		data := []interface{}{r.date, "ICOMN"}
		for _, i := range COMNR {
			data = append(data, s[i])
		}
		zt := "C"
		if data[6] == "0420" {
			zt = "Z"
		} else if s[9][0] == '0' || s[9][0] == '1' {
			zt = "D"
		}
		data = append(data, zt)
		stmt.Exec(data...)
		//fmt.Println(data...)
	}
	fmt.Println("导入 ICOMN 文件完成 ")
}

// LoadIERR 读取 IERR 文件
func (r *Reader) LoadIERR(path *zip.File, tx Execer) {
	offsets := []int{0, 3, 15, 27, 34, 45, 65, 78, 83, 90, 95, 104, 117, 120, 127, 139,
		151, 158, 161, 165, 178, 192, 204, 217, 230, 243, 248, 312, 323,
		327, 329, 331, 336, 349}
	columns := []int{0, 3, 4, 5, 6, 7, 8, 16, 30}
	data, _ := ReadAll(path, false)
	buf := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buf)

	Exec(tx, "delete from yl where lx=? and rq=?", "IERR", r.date)
	stmt, err := tx.Prepare("insert into yl(rq,lx,bz,seqno,cendt,kh,je,bwlx,jylx,oldseq,olddt,zt)values(?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic("准备 SQL 语句失败")
	}
	for scanner.Scan() {
		s := SplitData(scanner.Bytes(), offsets, columns)
		data := []interface{}{r.date, "IERR"}
		for _, i := range s {
			data = append(data, i)
		}
		zt := "C"
		if s[0] == "E23" {
			zt = "C"
		} else if s[6][0] == '0' || s[6][0] == '1' {
			zt = "D"
		}
		data = append(data, zt)
		stmt.Exec(data...)
		//fmt.Println(data...)
	}
	fmt.Println("导入 IEER 文件完成 ")
}

// LoadITFL 导入 ITFL 文件
func (r *Reader) LoadITFL(path *zip.File, tx Execer) {
	offsets := []int{0, 11, 18, 29, 49, 62, 67, 74, 79, 88, 101, 104, 111, 123, 130,
		133, 137, 150, 164, 176, 189, 201, 221, 233, 253, 257, 259, 261,
		263, 266, 277, 289, 291, 313}
	columns := []int{1, 2, 3, 4, 5, 6, 13, 29}
	data, _ := ReadAll(path, false)
	buf := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buf)

	Exec(tx, "delete from yl where lx=? and rq=?", "ITFL", r.date)
	stmt, err := tx.Prepare("insert into yl(rq,lx,seqno,cendt,kh,je,bwlx,jylx,oldseq,olddt,zt)values(?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic("准备 SQL 语句失败")
	}
	for scanner.Scan() {
		s := SplitData(scanner.Bytes(), offsets, columns)
		data := []interface{}{r.date, "ITFL"}
		for _, i := range s {
			data = append(data, i)
		}
		zt := "C"
		if s[5][0] == '0' || s[5][0] == '1' {
			zt = "D"
		}
		data = append(data, zt)
		stmt.Exec(data...)
	}
	fmt.Println("导入 ITFL 文件完成 ")
}

package xyk

import (
	"bufio"
	"bytes"
	"strings"
)

// LoadTrac 读取 TRAC 文件
func (r *Reader) LoadTrac(tx Execer) {
	data, _ := ReadAll(r.trac, true)
	buf := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.Contains(s, "440001") {
			fields := strings.Fields(s)
			Exec(tx, "insert or replace into ysqs values(?,?,?)", r.date, Atoi(fields[3]), Atoi(fields[4]))
			break
		}
	}
}

// LoadJorj 读取 JORJ文件
func (r *Reader) LoadJorj(tx Execer) {
	data, _ := ReadAll(r.jorj, true)
	buf := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buf)
	stmt, err := tx.Prepare("insert into jorj values(?,?,?,?,?)")
	defer stmt.Close()
	CheckErr(err)
	Exec(tx, "delete from jorj where rq=?", r.date) // 删除已导入的数据
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 7 && fields[2] == "440001" {
			stmt.Exec(r.date, fields[0], fields[1], Atoi(fields[5]), Atoi(fields[6]))
		}
	}
}
func extract(bytes []byte, start int, width int) string {
	return string(bytes[start-1 : start+width-1])
}

// LoadEve 读取 EVE 文件
func (r *Reader) LoadEve(tx Execer) {
	data, _ := ReadAll(r.eve, false)
	buf := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buf)
	Exec(tx, "delete from eve where rq=?", r.date) // 删除已导入的数据
	stmt, err := tx.Prepare("insert into eve values(?,?,?,?,?,?,?,?)")
	CheckErr(err)
	today := r.date[4:]
	prevday := PrevDay(r.date)
	for scanner.Scan() {
		row := scanner.Bytes()
		seqno := extract(row, 12, 6)
		cendt := extract(row, 18, 10)
		kh := strings.TrimSpace(extract(row, 28, 19))
		je := Atoi(extract(row, 47, 12))
		jdbz := extract(row, 59, 1)
		qsrq := extract(row, 113, 4)
		if qsrq == today {
			qsrq = r.date
		} else {
			qsrq = prevday
		}
		lsh := extract(row, 152, 6)
		stmt.Exec(r.date, seqno, cendt, kh, je, jdbz, qsrq, lsh)
	}
}

package xyk

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"strings"
)

func (r *Reader) LoadTrac(tx *sql.Tx) {
	data, _ := ReadAll(r.trac)
	buf := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.Contains(s, "440001") {
			fields := strings.Fields(s)
			tx.Exec("insert or replace into qsb(rq,ysjf,ysdf)values(?,?,?)", r.date, Atoi(fields[3]), Atoi(fields[4]))
			break
		}
	}
}

func (r *Reader) LoadRd1002(tx *sql.Tx) {
	data, _ := ReadAll(r.rd1002)
	buf := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 10 && fields[0] == "收付费" {
			_, err := tx.Exec("update qsb set sffjf+?,sffdf=? where rq=?", Atoi(fields[2]), Atoi(fields[3]), r.date)
			if err != nil {
				fmt.Println(err)
			}
		} else if len(fields) == 11 && fields[0] == "总" {
			_, err := tx.Exec("update qsb set jybs=?,jyjejf=?,jyjedf=?,jhfjf=?,jhfdf=?,qsfjf=?,qsfdf=?,qsjejf=?,qsjedf=? where rq=?",
				Atoi(fields[2]), Atoi(fields[3]), Atoi(fields[4]), Atoi(fields[5]), Atoi(fields[6]), Atoi(fields[7]), Atoi(fields[8]), Atoi(fields[9]), Atoi(fields[10]), r.date)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

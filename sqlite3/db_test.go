package sqlite3

import (
	"fmt"
	"testing"
)

func TestOpen(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Errorf("Open failed")
	}
	defer db.Close()
}

func TestExecMany(t *testing.T) {
	db, _ := Open(":memory:")
	defer db.Close()
	_, err := db.Exec(`create table if not exists abc(
		a		text    primary key,
		b 		text 	not null); --test
		`)
	if err != nil {
		t.Errorf("Execute Failed")
	}
	tr := db.Tran()
	tr.Add(`insert into abc values(?,?)`, "a", "b")
	tr.Add(`insert into abc values(?,?)`, "c", "b")
	tr.Exec()
	i := 0
	rd, _ := db.Fetch(`select * from abc`)
	for rd.Next() {
		fmt.Println(rd.Read()...)
		i += 1
	}
	var k int
	db.FetchValue(`select count(*) from abc`,&k)
	if k!=2{
		t.Errorf("Execute Failed")
	}
}


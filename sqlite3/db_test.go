package sqlite3

import (
	"testing"
)

func TestOpen(t *testing.T) {
	Config(":memory:")
	NewDB()
	defer Close()
}

func TestExecMany(t *testing.T) {
	Config(":memory:")
	defer Close()
	ExecScript(`create table if not exists abc(
		a		text    primary key,
		b 		text 	not null); --test
		insert into abc values(1,2);
		insert into abc values(3,4);
		`)
	val := FetchValue("select b from abc where a=?", 1)
	if val.(string) != "2" {
		t.Fatal("test FetchValue failed")
	}
	ExecTx(
		NewTr("insert into abc values(?,?)", 10, 20),
	)
	val = FetchValue("select b from abc where a=?", 10)
	if val.(string) != "20" {
		t.Fatal("test NewTr failed")
	}
	Println("select * from abc")
}

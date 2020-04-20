package sqlite3

import (
	"fmt"
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
	ExecScripts(`create table if not exists abc(
		a		text    primary key,
		b 		text 	not null); --test
		insert into abc values(1,2);
		insert into abc values(3,4);
		`)
	val := FetchValue("select b from abc where a=?", 1)
	if val.(string) != "2" {
		t.Fatal("test FetchValue failed")
	}
	r := Fetch("select * from abc")
	for r.Next() {
		fmt.Println(r.Read()...)
	}
}

package params

import (
	"grape/path"
	"grape/sqlite3"
	"testing"
)

func test() (err error) {
	tx := sqlite3.NewTx()
	defer tx.Rollback()
	path := path.NewPath("~/test.xlsx")
	err = LoadCheck(tx, "test", path, "1.0")
	tx.Commit()
	sqlite3.Println("select * from LoadFile")
	return
}
func TestLoader(t *testing.T) {
	if test() != nil {
		t.Error("Test Failed")
	}
	if test() == nil {
		t.Error("Test Failed")
	}
}

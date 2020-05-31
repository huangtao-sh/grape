package text

import (
	"bytes"
	"fmt"
	"testing"
)

func testConverter(t *testing.T) {
	bf := `1,2,3,4
5,6,7,8
9,10,11,12`
	b := bytes.NewReader([]byte(bf))
	r := NewReader(b, NewSepSpliter(","))
	for r.Next() {
		fmt.Println(r.Read()...)
	}
	t.Errorf("Failed")
}

package text

import (
	"bytes"
	"grape/data"
	"testing"
)

func TestConverter(t *testing.T) {
	bf := `1,2,hello,world
5,6,7,8
9,10,11,12`
	b := bytes.NewReader([]byte(bf))
	r := NewReader(b, false, NewSepSpliter(","), Exclude(0, 2), Include(1))
	d := data.NewData()
	d.Add(1)
	go d.Println()
	go r.ReadAll(d)
	d.Wait()
	//t.Errorf("hello")
}

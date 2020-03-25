package text

import "testing"

func testConverter(t *testing.T) {
	s := []interface{}{0, 1, 2, 3, 4, 5}
	c := NewIncluder(1, 4)
	d := c.Convert(s)
	if d[0].(int) != 1 && len(d) == 2 {
		t.Error("test Includer failed")
	}
	e := NewExcluder(0, 2)
	d = e.Convert(s)
	if d[0].(int) != 1 && len(d) == 3 {
		t.Error("test Exclude failed")
	}
}

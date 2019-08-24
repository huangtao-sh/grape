package path

import (
	"testing"
)

func TestExpand(t *testing.T) {
	p := NewPath("~hunter/abc")
	if p.String() != `C:\Users\hunter\abc` {
		t.Errorf("Test Failed!")
	}

	path := Expand("$programdata/abc")
	if path != `C:\ProgramData\abc` {
		t.Errorf("Test Failed!")
	}
	path = Expand("%programdata%/abc")

	if path != `C:\ProgramData\abc` {
		t.Errorf("Test Expand Failed!")
	}
}

func TestTempDir(t *testing.T) {
	if TempDir != `C:\Users\huangtao\AppData\Local\Temp` {
		t.Errorf("Test TempDir Failed")
	}
}

package csvfile

import (
	"encoding/csv"
	"grape/path"
	"os"
	"strings"
	"testing"
	"grape/gbk"
	//"io/ioutil"
)

func createTmp() (filename string, err error) {
	filename = path.TempDir + "/abc.csv"
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer f.Close()
	w := csv.NewWriter(gbk.NewWriter(f))
	for i := 0; i < 3; i++ {
		w.Write([]string{"张三", "b"})
	}
	w.Flush()
	return
}

func TestReader(t *testing.T) {
	filename, err := createTmp()

	if err != nil {
		t.Errorf("创建文件失败")
	}
	defer os.Remove(filename)

	r, err := NewReader(filename, "GBK")
	if err != nil {
		t.Errorf("打开文件失败")
	}
	defer r.Close()
	i := 0
	for {
		record, err := r.Read()
		if err != nil {
			break
		}
		s, _ := record[0].(string)
		if s != "张三" {
			t.Errorf("Test Read1 Failed")
		}
		i++
	}
	if i != 3 {
		t.Errorf("Test Read Failed")
	}
}

var data = `abc,def
abc,def
abc,def`

func TestReader2(t *testing.T) {
	r := csv.NewReader(strings.NewReader(data))
	i := 0
	for {
		record, err := r.Read()
		if err != nil {
			break
		}
		if record[0] != `abc` {
			t.Errorf("Test Read1 Failed")
		}
		i++
	}
	if i != 3 {
		t.Errorf("Test Read Failed")
	}
}

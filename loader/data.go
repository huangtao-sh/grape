package loader

// Reader 数据读取接口
type Reader interface {
	Read() ([]string, error)
}

// ConvertFunc 数据转换类型
type ConvertFunc func([]string) ([]string, error)

// Converter 数据转换类
type Converter struct {
	Reader     // 源数据读取
	converters []ConvertFunc
}

// NewConverter Converter 构造函数
func NewConverter(r Reader, converters ...ConvertFunc) *Converter {
	return &Converter{r, converters}
}

// Read 读取方法
func (c *Converter) Read() (rows []string, err error) {
	rows, err = c.Reader.Read()
	for _, conv := range c.converters {
		if err != nil  || rows == nil {
			return
		}
		rows, err = conv(rows)
	}
	return
}

// Slice 把 []string 转换成 []interface{}
func Slice(row []string) (col []interface{}) {
	col = make([]interface{}, len(row))
	for i, v := range row {
		col[i] = v
	}
	return
}

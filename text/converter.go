package text

// Includer 转换函数，仅保留指定行，实现行转换接口
type Includer struct {
	columns []int
}

// NewIncluder 包含列构造函数
func NewIncluder(columns ...int) *Includer {
	return &Includer{columns}
}

// Convert 转换函数
func (i *Includer) Convert(s []interface{}) (d []interface{}) {
	columns := i.columns
	d = make([]interface{}, len(columns))
	for i, idx := range columns {
		d[i] = s[idx]
	}
	return
}

// Excluder 排除
type Excluder struct {
	columns map[int]bool
}

// NewExcluder 包含列构造函数
func NewExcluder(columns ...int) *Excluder {
	column := make(map[int]bool)
	for _, i := range columns {
		column[i] = true
	}
	return &Excluder{column}
}

// Convert 转换函数
func (i *Excluder) Convert(s []interface{}) (d []interface{}) {
	columns := i.columns
	for i, value := range s {
		if !columns[i] {
			d = append(d, value)
		}
	}
	return
}

// Convert 转换数据函数模型
type Convert func(interface{}) interface{}

// Converter 行转换
type Converter struct {
	converter map[int]Convert
}

// NewConverter 行转换器构造函数
func NewConverter(converter map[int]Convert) *Converter {
	return &Converter{converter}
}

// Convert 转换
func (c *Converter) Convert(s []interface{}) []interface{} {
	for idx, convert := range c.converter {
		s[idx] = convert(s[idx])
	}
	return s
}

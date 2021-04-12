package sqlite3

import (
	"database/sql"
	"errors"
	"fmt"
	"grape/path"
	"grape/util"
	"strings"

	// 引用 go-sqlite3 进行初始化
	_ "github.com/mattn/go-sqlite3"
)

var (
	dataSourceName string  // 默认数据连接
	db             *sql.DB // 数据连接
)

// parsePath 数据库配置处理函数
func parsePath(pathName string) string {
	// path 如果不是 :memory:，无目录的指定默认目录
	if (pathName != ":memory:") && (path.NewPath(pathName).Dir() == ".") {
		dataHome := path.Home.Join(".data")
		dataHome.Ensure() // 目录不存在则自动创建
		pathName = (dataHome.Join(pathName).WithExt(".db")).String()
	}
	return pathName
}

// Config  配置数据库文件
// 配置完成之后，可以直接调用 Open 打开
func Config(pathName string) {
	Close() // 先关闭当前的数据库连接
	dataSourceName = parsePath(pathName)
}

// NewDB 打开数据库
func NewDB() *sql.DB {
	if db == nil {
		var err error
		db, err = sql.Open("sqlite3", dataSourceName)
		util.CheckFatal(err)
	}
	return db
}

// Attach 附加数据库
func Attach(path string, name string) {
	_, err := NewDB().Exec(fmt.Sprintf("attach database '%s' as '%s'", parsePath(path), name))
	util.CheckFatal(err)
}

// Detach 分离数据库
func Detach(name string) {
	_, err := NewDB().Exec(fmt.Sprintf("detach database '%s'", name))
	util.CheckFatal(err)
}

// Close 关闭数据库
func Close() {
	if db != nil { // 判断当前连接是否有值，如有则关闭
		db.Close()
		db = nil
	}
}

// ExecScript 执行脚本，支持多条语句，用于执行 DDL 语句
func ExecScript(sql string) {
	_, err := NewDB().Exec(sql)
	util.CheckFatal(err)
}

// LoadSQL 生成导入 SQL 语句
func LoadSQL(method string, table string, columns interface{}) string {
	var colCount int
	switch cols := columns.(type) {
	case string:
		fields := strings.Split(cols, ",")
		colCount = len(fields)
		table = fmt.Sprintf("%s(%s)", table, cols)
	case int:
		colCount = cols
	case []string:
		colCount = len(cols)
		table = fmt.Sprintf("%s(%s)", table, strings.Join(cols, ","))
	default:
		util.CheckFatal(errors.New("columns must be string,int or []stirng"))
	}
	s := fmt.Sprintf("%%s into %%s %%%dV", colCount)
	return util.Sprintf(s, method, table)
}

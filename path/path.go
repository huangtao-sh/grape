package path

import (
	"fmt"
	"grape/date"
	"grape/util"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var (
	// HomeDir 用户家目录
	HomeDir string
	// TempDir 临时目录
	TempDir string
	// Home 家目录
	Home *Path
)

func init() {
	ur, _ := user.Current()
	HomeDir = ur.HomeDir
	TempDir = os.TempDir()
	Home = NewPath("~")
}

// Expand 扩展路径
func Expand(p string) (path string) {
	var root bool
	if len(p) > 0 && p[0] == '/' {
		root = true
	}
	parts := strings.Split(filepath.ToSlash(p), "/")
	if len(parts) >= 1 {
		base := parts[0]
		if strings.HasPrefix(base, "~") {
			parts[0] = HomeDir
			if len(base) > 1 {
				parts[0] = parts[0] + "/../" + base[1:]
			}
		} else if matched, _ := filepath.Match("%*%", base); matched {
			parts[0] = "$" + strings.Replace(base, "%", "", -1)
		}
		if len(base) == 2 && base[1] == ':' {
			parts[0] = base + "/"
		}
	}
	p = filepath.Join(parts...)
	path = os.ExpandEnv(p)
	if root {
		path = fmt.Sprintf("/%s", path)
	}
	return
}

// Path 目录
type Path struct {
	Raw, path string
}

// NewPath 目录构造函数
func NewPath(path string) (p *Path) {
	return &Path{path, Expand(path)}
}

// Remove 删除文件
func (p *Path) Remove() error {
	return os.Remove(p.path)
}

// Open 打开文件
func (p *Path) Open() (io.ReadCloser, error) {
	return os.Open(p.path)
}

// String 转换成字符串
func (p *Path) String() string {
	return p.path
}

// IsExist 判断目录是否存在
func (p *Path) IsExist() bool {
	_, err := os.Stat(p.path)
	return err == nil || os.IsExist(err)
}

// MakeDir 创建目录
func (p *Path) MakeDir() error {
	return os.Mkdir(p.path, 0755) // 调整生成文件的权限为 755
}

// MakeDirAll 创建目录
func (p *Path) MakeDirAll() error {
	return os.MkdirAll(p.path, 0755) // 调整生成文件的权限为 755
}

// Ensure 确保目录存在
func (p *Path) Ensure() error {
	if !p.IsExist() {
		return p.MakeDirAll()
	}
	return nil
}

// Dir 目录
func (p *Path) Dir() string {
	return filepath.Dir(p.path)
}

// Ext 扩展名
func (p *Path) Ext() string {
	return filepath.Ext(p.path)
}

// Base 文件名
func (p *Path) Base() string {
	return filepath.Base(p.path)
}

// Join 连接目录
func (p *Path) Join(elem ...string) (path *Path) {
	pa := filepath.Join(p.path, filepath.Join(elem...))
	return NewPath(pa)
}

// Clean 整理文件
func (p *Path) Clean() string {
	return filepath.Clean(p.path)
}

// Glob 查找文件
func (p *Path) Glob(pattern string) (matches []string) {
	matches, err := filepath.Glob(filepath.Join(p.path, pattern))
	util.CheckFatal(err)
	return
}

// Match 判断是否匹配
func (p *Path) Match(pattern string) (matched bool) {
	matched, err := filepath.Match(pattern, p.path)
	util.CheckFatal(err)
	return
}

// Find 在目录中查找记录
func (p *Path) Find(pattern string) (path string) {
	matches := p.Glob(pattern)
	if len(matches) > 0 {
		path = matches[len(matches)-1]
	}
	return
}

// FileInfo 获取文件信息
func (p *Path) FileInfo() (info os.FileInfo) {
	info, err := os.Stat(p.String())
	util.CheckFatal(err)
	return
}

// Split 拆分文件名，将文件名拆分行目录和文件名
func (p *Path) Split() (dir, name string) {
	return filepath.Split(p.path)
}

// WithExt 替换成指定的扩展名
func (p *Path) WithExt(ext string) *Path {
	dir, name := p.Split()
	name = strings.Split(name, ".")[0]
	return NewPath(dir).Join(name + ext)
}

// InitLog 初始化日志文件
// 放在 $HOME/.logs/date 目录下，并且用命令作为文件名
func InitLog() {
	_, name := NewPath(os.Args[0]).Split() // 获取命令行
	dir := NewPath("~/.logs").Join(date.Today().Format("%F"))
	dir.Ensure() // 建立目录
	filename := dir.Join(name).WithExt(".log")
	r, err := os.OpenFile(filename.String(), os.O_CREATE|os.O_APPEND, os.ModePerm)
	util.CheckFatal(err)
	log.SetOutput(r)
}

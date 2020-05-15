package path

import (
	"grape/util"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// HomeDir 用户家目录
var HomeDir string

// TempDir 临时目录
var TempDir string

// Home 家目录
var Home *Path

func init() {
	ur, _ := user.Current()
	HomeDir = ur.HomeDir
	TempDir = os.TempDir()
	Home = NewPath("~")
}

// Expand 扩展路径
func Expand(p string) (path string) {
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
func (p *Path) Open() (*os.File, error) {
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
	return os.Mkdir(p.path, os.ModeDir)
}

// MakeDirAll 创建目录
func (p *Path) MakeDirAll() error {
	return os.MkdirAll(p.path, os.ModeDir)
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

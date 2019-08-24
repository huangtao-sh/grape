package path

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"


)

var HomeDir, TempDir string

func init() {
	ur, _ := user.Current()
	HomeDir = ur.HomeDir
	TempDir = os.TempDir()
}

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
	}
	p = filepath.Join(parts...)
	path = os.ExpandEnv(p)
	return
}


type Path struct {
	Raw, path string
}

func NewPath(path string) (p *Path) {
	return &Path{path, Expand(path)}
}

// 删除文件
func (p *Path) Remove() error {
	return os.Remove(p.path)
}

// 打开文件
func (p *Path) Open() (*os.File, error) {
	return os.Open(p.path)
}

// 转换成字符串
func (p *Path) String() string {
	return p.path
}

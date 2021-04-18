package util

import (
	"grape"
	"grape/date"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// InitLog 初始化 log
func InitLog() {
	cmd := filepath.Base(os.Args[0])
	name := strings.Split(cmd, ".")[0] + ".log"
	home, _ := os.UserHomeDir()
	logRoot := filepath.Join(home, ".logs", date.Today().String())
	os.Mkdir(logRoot, 0755)
	filePath := filepath.Join(logRoot, name)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	grape.CheckFatal(err)
	log.SetOutput(f)
}

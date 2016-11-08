// 日志处理
package illusion

import (
	"os"
	"strings"
	"time"
	//	"sync"
	//"fmt"
	"log"
	"path/filepath"
	//"sync"
)

//通过设置返回对应的logger对象
func setLogger(basePath string) {
	handler := findTodayLogFile(basePath)
	forLogger(handler)
}

func findTodayLogFile(relativePath string) *os.File {
	basePath := string(filepath.Separator) + strings.TrimPrefix(relativePath, "/")
	absPath, _ := filepath.Abs(".")
	fullPathToTodayFile := absPath + basePath + string(filepath.Separator) + strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")[0] + ".log"
	handler, err := os.OpenFile(fullPathToTodayFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil { //直接让程序崩溃
		panic("创建日志文件失败:" + err.Error())
	}
	return handler
}

//设置相对应的logger设置
func forLogger(file *os.File) {
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

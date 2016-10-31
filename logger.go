// 日志处理
package illusion

import (
	"os"
	"time"
	"strings"
//	"sync"
	"path/filepath"
	"log"
	"sync"
	"fmt"
)
type ILogger interface {
	//这一个接口就好了
	log(string) error
}


//暂时设置这个logger为单例
//提供两个可调用的接口来设置
type Logger struct {
	//logger基础路径
	basePath  string

	//持有的文件句柄
	file    *os.File

	//错误信息
	Error error
}

var logger *Logger
var once sync.Once
//var mutex sync.Mutex
//var logPath string

//每个程序启动的时候只能被初始化一次
func setLogger(path string){
	once.Do(func(){
		logger = outNewLogger(path)
	})
}

//全局单例
func loggerInstance()*Logger{
	return logger
}

//通过设置返回对应的logger对象
func outNewLogger(basePath string)*Logger{
	l := &Logger{basePath: basePath, Error:nil}
	writer := l.findTodayFile()
	l.setLogger(writer)
	return l
}

//找到logger目录基于文件系统的目录
func (it *Logger)systemDir() string{
	basePath := string(filepath.Separator) + strings.TrimPrefix(it.basePath, "/")
	absPath,err := filepath.Abs(".")
	if err != nil {
		it.Error = err
		return ""
	}
	dir := absPath + basePath + string(filepath.Separator)
	return dir
}

//找到今天的logger文件，没有则创建
//可以肯定的是这里的log file可以抽象为任意的目的地，可以是缓存，数据库，当然也包括文件
func (it *Logger)findTodayFile()*os.File{
	todayFile := strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")[0] + ".log"
	logDir := it.systemDir()
	if it.Error != nil {
		return nil
	}
	fullTodayFilePath := logDir + todayFile
	//today := split[0] + ".log"
	handler,err := os.OpenFile(fullTodayFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		it.Error = err
		return nil
	}
	return handler
}

//设置相对应的logger设置
func (it *Logger)setLogger(file *os.File){
	if it.Error != nil {
		return
	}
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

//打印log文件
func (it *Logger)Log(msg string){
	if it==nil {
		fmt.Println("怎么可能")
		return
	}
	if it.Error != nil {
		//这好像是一个悖论
		log.Println(it.Error.Error())
		return
	}
	log.Println(msg)
}
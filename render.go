//模板设置
package illusion

import (
	//"html/template"
	"strings"
	//"io"
	"bytes"
	"errors"
	//"fmt"
	"github.com/flosch/pongo2"
)

type Value pongo2.Context

//好像只要一个就行了
//所以Context只需要持有这个接口就好了
type IllusionTemplate interface {
	//渲染文件
	//获取文件内容
	Content(string, interface{}) []byte

	//清理Buffer中原有的内容
	//Clear()
	//这个起个什么名字好呢
	//这个名字好像很不错
	//echo(...interface{})

	//获取错误信息
	//Error() error
}

//存储模板引擎之用
type Template struct {
	//基础view路径
	baseFileLocation string

	//具体加载到的template文件
	//template *template.Template
	//这个好像没卵用

	//同时还需要一个专门用以获取模板内容的Writer
	contentWriter *ContentWriter
	//错误信息
	//Err error
}

//每个Context会附着一个Template
//illusion中处理basePath
//为了避免每次都需要在这里重新拼接字符串
func newTemplate(basePath string, writer *ContentWriter) *Template {
	return &Template{baseFileLocation: basePath, contentWriter: writer}
}

//渲染这个文件
//返回最终结果 string
func (it *Template) Content(file string, value Value) []byte {
	file = strings.TrimPrefix(file, "/")
	finalPath := it.baseFileLocation + file
	tpl,err := pongo2.FromFile(finalPath)
	if err != nil {
		appendError(errorInfo{Error:err, Level:panicOnError})
	}
	result,err := tpl.Execute(value)
	if err != nil {
		appendError(errorInfo{Error:err, Level:panicOnError})
	}
	return []byte(result)
	/*t, err := template.ParseFiles(finalPath)
	if err != nil {

		return
	}
	err = t.Execute(it.contentWriter, value)
	if err != nil {
		appendError(errorInfo{Error:err, Level:panicOnError})
	}
	result = it.contentWriter.Read()
	return*/
}

//哎呦，层层调用
//:)
func (it *Template) Clear() {
	it.contentWriter.Clear()
}

type ContentWriter struct {
	buf   bytes.Buffer
	//Error error
}

func newContentWriter() *ContentWriter {
	return &ContentWriter{buf: bytes.Buffer{}}
}

func (writer *ContentWriter) Write(p []byte) (n int, err error) {
	n, err = writer.buf.Write(p)
	if n < len(p) {
		appendError(errorInfo{Error:errors.New("未完全获取模板内容"), Level: logOnError})
		//writer.Error = errors.New("没写全")
	}
	return n, err
}

//获取写入的内容
//直接返回字符数组就好
func (writer *ContentWriter) Read() []byte {
	//if writer.Error != nil {
	//	return nil
	//}
	return writer.buf.Bytes()
}

//清理缓存
func (writer *ContentWriter) Clear() {
	writer.buf.Reset()
}

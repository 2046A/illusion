//模板设置
package illusion

import (
	"html/template"
	"strings"
	//"io"
	"bytes"
	"errors"
	//"fmt"
)

//好像只要一个就行了
//所以Context只需要持有这个接口就好了
type IllusionTemplate interface {
	//渲染文件
	//获取文件内容
	Content(string, interface{}) []byte

	//清理Buffer中原有的内容
	Clear()
	//这个起个什么名字好呢
	//这个名字好像很不错
	//echo(...interface{})

	//获取错误信息
	Error() error
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
	Err error
}

//每个Context会附着一个Template
//illusion中处理basePath
//为了避免每次都需要在这里重新拼接字符串
func newTemplate(basePath string, writer *ContentWriter) *Template {
	return &Template{baseFileLocation: basePath, contentWriter: writer}
}

//渲染这个文件
//返回最终结果 string
func (it *Template) Content(file string, value interface{}) (result []byte) {
	file = strings.TrimPrefix(file, "/")
	finalPath := it.baseFileLocation + file
	t, err := template.ParseFiles(finalPath)
	if err != nil {
		it.Err = err
		//fmt.Println("文件:" + finalPath)
		//result = []byte("")
		return
	}
	//fmt.Println("serving file:" + finalPath)
	//contentWriter := newContentWriter()
	t.Execute(it.contentWriter, value)
	//if it.contentWriter.Error != nil {
	//panic("出错了，未读出所有的内容")
	//}
	//if t.contentWriter
	result = it.contentWriter.Read()
	//fmt.Println("*****************************************************")
	//fmt.Println(string(result))
	//fmt.Print("*****************************************************")
	return
}

//哎呦，层层调用
//:)
func (it *Template) Clear() {
	it.contentWriter.Clear()
}

//返回相应的错误信息
func (it *Template) Error() error {
	return it.contentWriter.Error
}

type ContentWriter struct {
	buf   bytes.Buffer
	Error error
}

func newContentWriter() *ContentWriter {
	return &ContentWriter{buf: bytes.Buffer{}}
}

func (writer *ContentWriter) Write(p []byte) (n int, err error) {
	n, err = writer.buf.Write(p)
	if n < len(p) {
		writer.Error = errors.New("没写全")
	}
	return n, err
}

//获取写入的内容
//直接返回字符数组就好
func (writer *ContentWriter) Read() []byte {
	if writer.Error != nil {
		return nil
	}
	return writer.buf.Bytes()
}

//清理缓存
func (writer *ContentWriter) Clear() {
	writer.buf.Reset()
}

//我就实现一个最简单的context好了
//提高可读性

package illusion

import (
	"net/http"
	//"go/token"
	"net"
	"strings"
	"time"
)

//这个好像没什么用
const abortIndex int8 = 10
const MaxParamSize = 20 //最大url参数个数，是指/user/:id中id这样的个数
//const

//Context还是先简单着来
//因为我还闹不大明白很多东西, :)
type Context struct {
	//请求对象
	Request  *http.Request

	//response对象
	Writer   http.ResponseWriter

	//额外参数
	Params   Params

	//调用链
	handlers HandlerChain

	//额外附着在Context上的数据
	Keys     map[string]interface{}

	//错误信息
	Error    error

	//是否需要终止
	aborted  bool
}

//初始化一个Context
func newContext() *Context {
	return &Context{
		Request:  nil,
		Writer:   nil,
		Params:   make(Params, 0, MaxParamSize),
		handlers: make(HandlerChain, 0, MaxHandlerNumber),
		Keys:     make(map[string]interface{}),
		Error:    nil,
		aborted: false,
	}
}

//Context是从pool获取的，使用前必须调用reset清理原来的值
func (it *Context) reset() {
	it.Params = it.Params[0:0]     //这干嘛的
	it.handlers = it.handlers[0:0] //..
	it.Request = nil
	it.Writer = nil
	it.Keys = make(map[string]interface{})
	it.Error = nil
	it.aborted = false
}

//好像还有个问题, 让我想想？？？？？？？
//abort的时候要确保能终止
//让我想想 ....
//还在想
//还在想..
func (it *Context) AbortWithStatus(code int) {
	it.Status(code)
	it.Abort()
}

func (it *Context) AbortWithError(code int, err error) {
	it.Error = err
	it.AbortWithStatus(code)
}

//设置结束标志
func (it *Context) Abort() {
	it.aborted = true
}

//判断是否结束
//好像没什么用
func (it *Context) IsAborted() bool {
	return it.aborted
}

//开始处理http请求
//这是整个request的最终处理地点
func (it *Context) Next() {
	for _, handler := range it.handlers {
		//那么问题来了
		// 整个调用链中的其中一个函数如何使整个链结束呢
		//在context中添加Abort标志以示结束
		handler(it)
		//如果用户设置了Abort, 那么结束
		if it.aborted {
			return
		}
	}
}

//添加额外的信息到调用链中
//后面可以使用
func (it *Context) Append(key string, value interface{}) {
	if it.Keys == nil {
		it.Keys = make(map[string]interface{})
	}
	it.Keys[key] = value
}

//获取附加的信息
func (it *Context) Retrieve(key string) (value interface{}, exists bool) {
	if it.Keys != nil {
		value, exists = it.Keys[key]
	}
	return
}

/*******************************
/*** request读数据 *************
/******************************/

//比如:/user/:id ,那么it.Param("id")就会返回相应的:id值
func (it *Context) Param(key string) string {
	return it.Params.ByName(key)
}

//是it.Request.URL.Query()的快捷方式
// GET /?name=Manu&lastname=
// ("Manu", true) == it.GetQuery("name")
// ("", false) == it.GetQuery("id")
// ("", true) == it.GetQuery("lastname")
func (it *Context) GetQuery(key string) (string, bool) {
	req := it.Request
	if values, ok := req.URL.Query()[key]; ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}

//诸如/path?id=1234&name=Manu
//it.Query("id") == 1234
//it.Query("name") == "Manu"
func (it *Context) Query(key string) string {
	value, _ := it.GetQuery(key)
	return value
}

//获取参数的默认版本
func (it *Context) DefaultQuery(key, defaultValue string) string {
	if value, ok := it.GetQuery(key); ok {
		return value
	}
	return defaultValue
}

//与getQuery行为一致，不做过多解释
func (it *Context) GetPostForm(key string) (string, bool) {
	req := it.Request
	req.ParseMultipartForm(32 << 20) //32MB ???
	if values := req.PostForm[key]; len(values) > 0 {
		return values[0], true
	}
	if req.MultipartForm != nil && req.MultipartForm.File != nil {
		if values := req.MultipartForm.Value[key]; len(values) > 0 {
			return values[0], true
		}
	}
	return "", false
}

//如果key存在，返回对应值，否则返回""
func (it *Context) PostForm(key string) string {
	value, _ := it.GetPostForm(key)
	return value
}

func (it *Context) DefaultPostForm(key, defaultValue string) string {
	if value, ok := it.GetPostForm(key); ok {
		return value
	}
	return defaultValue
}

//还有一些简单的东西可以实现以下

func (it *Context) requestHeader(key string) string {
	if values, _ := it.Request.Header[key]; len(values) > 0 {
		return values[0]
	}
	return ""
}

// ClientIP implements a best effort algorithm to return the real client IP, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
func (it *Context) ClientIP() string {
	//if c.illusion.ForwardedByClientIP {
	clientIP := strings.TrimSpace(it.requestHeader("X-Real-Ip"))
	if len(clientIP) > 0 {
		return clientIP
	}
	clientIP = it.requestHeader("X-Forwarded-For")
	if index := strings.IndexByte(clientIP, ','); index >= 0 {
		clientIP = clientIP[0:index]
	}
	clientIP = strings.TrimSpace(clientIP)
	if len(clientIP) > 0 {
		return clientIP
	}
	//}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(it.Request.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

//?????
func (it *Context) ContentType() string {
	//???????
	return it.requestHeader("Content-Type")
}

/*****************************
/*** Response写操作 **********
/****************************/

//写key->value到Header中
//如果value=""，也没差啊
func (it *Context) Header(key, value string) {
	if len(value) == 0 {
		it.Writer.Header().Del(key)
	} else {
		it.Writer.Header().Set(key, value)
	}
}

//被带到坑里去了
func (it *Context) Status(code int) {
	it.Writer.WriteHeader(code)
}

//这个Write应该只能被调用一次就好
//否则会报header被重复写的错误
func (it *Context) Write(status int, value string) {
	//这里出现了二次写头部的问题?????
	it.Status(status)
	it.Writer.Write([]byte(value))
}

//其他的一些再说，反正我也不懂

/*******************
 *** context ******
 ** 依我观察, 好像下面的没什么作用啊 ***********
 *******************/

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *Context) Done() <-chan struct{} {
	return nil
}

func (c *Context) Err() error {
	return nil
}

func (c *Context) Value(key interface{}) interface{} {
	if key == 0 {
		return c.Request
	}
	if keyAsString, ok := key.(string); ok {
		val, _ := c.Retrieve(keyAsString)
		return val
	}
	return nil
}

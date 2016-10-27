//可以参考python flask的blueprint功能,大致的功能是这样的
// backend := &Blueprint{BasePath:"/backend", Name:"backend"}
// backend.Before(HandlerFunc)  这里可以做验证，log等等操作
// backend.Get("/home", HandlerFunc) Url处理
// backend.After(HandlerFunc)   之后的操作????
// globalApp.register(backend)
// globalApp.Run(":8080")

package illusion

import (
	"regexp"
)

//每个httpMethod下面会有很多path->HandlerFunc的映射
type methodMap  map[string]HandlerFunc

//有什么很好的存储格式???
//这个存储格式太丑了 httpMethod -> MethodMap
type methodHashMap  map[string]methodMap

//Blueprint应该实现的接口
//供参考
type BluePrinter interface {
	//进入handler之前的操作
	Before(HandlerFunc)

	//进入handler之后的操作,?
	After(HandlerFunc)

	//返回自身，构成链式调用
	Handle(string, string, HandlerFunc) BluePrinter

	//匹配任意的http方法
	Any(string, HandlerFunc) BluePrinter

	//如下为标准的http方法
	Get(string, HandlerFunc) BluePrinter
	Post(string, HandlerFunc) BluePrinter
	Delete(string, HandlerFunc) BluePrinter
	Patch(string, HandlerFunc) BluePrinter
	Put(string, HandlerFunc) BluePrinter
	Options(string, HandlerFunc) BluePrinter
	Head(string, HandlerFunc) BluePrinter

	//这怎么用呢？
	//等我了解一下
	//我觉着还是交给核心路由器好了
	//StaticFile(string, string) BluePrinter
	//Static(string, string) BluePrinter
	//StaticFs(string, string) BluePrinter
}

//在Blueprint被注册到Illusion前，会持有所有的相关信息
type Blueprint struct {
	//对应此Blueprint的基础路径
	BasePath string

	//名称
	//名称应该唯一
	Name string

	//处理链, 可以一次性存储，
	//毕竟，你不会放太多的函数在调用链中, right? :)
	//这个好像也不需要了
	//Handlers HandlerChain

	//handler之前的处理链
	BeforeChain HandlerChain

	//handler之后的处理链
	AfterChain HandlerChain

	//错误信息
	Err error

	//必须要持有的核心路由
	//不显示持有也可以，但不明显，最好持有
	//核心路由在全局是个单例
	//illusion *Illusion
	//不用了

	//存储所有的相关信息
	MethodHashMap methodHashMap
}

//将handler合并进一个数组中
//func (it *Blueprint)extendBeforeChain(handler HandlerFunc){
//	return
//}

//func (it *Blueprint)extendAfterChain(Handler HandlerFunc){
///return
//}

//如下是Blueprint需要实现的接口
//func NewBlueprint() *Blueprint {
//	return Blueprint("/", "home")
//}

func Blueprint(path, name string) *Blueprint {
	return &Blueprint{
		BasePath:    path,
		Name:        name,
		BeforeChain: make(HandlerChain, 0, 5), //最大的beforeChain个数
		AfterChain:  make(HandlerChain, 0, 5), //最大的AfterChain个数
		MethodHashMap:  make(methodHashMap), //这个反倒好点
		Err:         nil,
		//illusion:    globalIllusion(),
	}
}

//结合beforeChain + Handler + afterChain形成一个调用链
func (it *Blueprint) fullChain(handler HandlerFunc) HandlerChain {
	if it.Err != nil {return }
	return nil
}

//结合blueprint.BasePath + relativePath形成绝对Url
func (it *Blueprint) truePath(relativePath string) string {
	if it.Err != nil {return }
	return CleanPath(it.BasePath + "/" + relativePath)
}

//我很怀疑append这个操作..
func (it *Blueprint) Before(handler HandlerFunc) *Blueprint {
	if it.Err != nil {return }
	it.BeforeChain = append(it.BeforeChain, handler)
	return it
}

//同上
func (it *Blueprint) After(handler HandlerFunc) *Blueprint {
	if it.Err != nil {return }
	it.AfterChain = append(it.AfterChain, handler)
	return it
}

func (it *Blueprint)handle(httpMethod, relativePath string, handler HandlerFunc) *Blueprint{
	finalPath := it.truePath(relativePath)
	it.MethodHashMap[httpMethod][finalPath] = handler
	return it
}

func (it *Blueprint) Handle(httpMethod, relativePath string, handler HandlerFunc) *Blueprint {
	if it.Err != nil {return}
	if matches, err := regexp.MatchString("^[A-Z]+$", httpMethod); !matches || err != nil {
		it.Err = err
		panic("http method " + httpMethod + " is not valid")
	}
	return it.handle(httpMethod, relativePath, handler)
}

func (it *Blueprint)Post(relativePath string, handler HandlerFunc) *Blueprint{
	return it.handle("POST", relativePath, handler)
}

func (it *Blueprint)Get(relativePath string, handler HandlerFunc) *Blueprint{
	return it.handle("GET", relativePath, handler)
}

func (it *Blueprint)Delete(relativePath string, handler HandlerFunc) *Blueprint{
	return it.handle("DELETE", relativePath, handler)
}

func (it *Blueprint)Patch(relativePath string, handler HandlerFunc) *Blueprint{
	return it.handle("PATCH", relativePath, handler)
}

func (it *Blueprint)Put(relativePath string, handler HandlerFunc) *Blueprint{
	return it.handle("PUT", relativePath, handler)
}

func (it *Blueprint)Options(relativePath string, handler HandlerFunc) *Blueprint{
	return it.handle("OPTIONS", relativePath, handler)
}

func (it *Blueprint)Head(relativePath string, handler HandlerFunc) *Blueprint{
	return it.handle("HEAD", relativePath, handler)
}

//...
func (it *Blueprint)Any(relativePath string, handler HandlerFunc) *Blueprint{
	//return it.handle("DELETE", relativePath, handler)
	it.handle("GET", relativePath, handler)
	it.handle("POST", relativePath, handler)
	it.handle("PUT", relativePath, handler)
	it.handle("PATCH", relativePath, handler)
	it.handle("HEAD", relativePath, handler)
	it.handle("OPTIONS", relativePath, handler)
	it.handle("DELETE", relativePath, handler)

	//这两个是什么???
	it.handle("CONNECT", relativePath, handler)
	it.handle("TRACE", relativePath, handler)
	return it
}


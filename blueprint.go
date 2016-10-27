//可以参考python flask的blueprint功能,大致的功能是这样的
// backend := &Blueprint{BasePath:"/backend", Name:"backend"}
// backend.Before(HandlerFunc)  这里可以做验证，log等等操作
// backend.Get("/home", HandlerFunc) Url处理
// backend.After(HandlerFunc)   之后的操作????
// globalApp.register(backend)
// globalApp.Run(":8080")

package illusion

import "regexp"

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
	StaticFile(string, string) BluePrinter
	Static(string, string) BluePrinter
	StaticFs(string, string) BluePrinter
}

type Blueprint struct {
	//对应此Blueprint的基础路径
	BasePath string

	//名称
	//名称应该唯一
	Name string

	//处理链, 可以一次性存储，
	//毕竟，你不会放太多的函数在调用链中,right? :)
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
	illusion *Illusion
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
		Err:         nil,
		illusion:    globalIllusion(),
	}
}

//结合beforeChain + Handler + afterChain形成一个调用链
func (it *Blueprint) fullChain(handler HandlerFunc) HandlerChain {
	return nil
}

//结合blueprint.BasePath + relativePath形成绝对Url
func (it *Blueprint) TruePath(relativePath string) string {
	return CleanPath(it.BasePath + "/" + relativePath)
}

func (it *Blueprint) Before(handler HandlerFunc) BluePrinter {
	it.BeforeChain = append(it.BeforeChain, handler)
	return it
}

func (it *Blueprint) After(handler HandlerFunc) BluePrinter {
	it.AfterChain = append(it.AfterChain, handler)
	return it
}

func (it *Blueprint) Handle(httpMethod, relativePath string, handler HandlerFunc) BluePrinter {
	if matches, err := regexp.MatchString("^[A-Z]+$", httpMethod); !matches || err != nil {
		it.Err = err
		panic("http method " + httpMethod + " is not valid")
	}
	handlers := it.fullChain(handler)
	finalPath := it.TruePath(relativePath)
	//这是illusion所缺少的啦
	it.illusion.addRoute(httpMethod, finalPath, handlers)
	return it
}


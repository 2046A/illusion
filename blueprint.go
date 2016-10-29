//可以参考python flask的blueprint功能,大致的功能是这样的
// backend := &Blueprint{BasePath:"/backend", Name:"backend"}
// backend.Before(HandlerFunc)  这里可以做验证，log等等操作
// backend.Get("/home", HandlerFunc) Url处理
// backend.After(HandlerFunc)   之后的操作????
// globalApp.register(backend)
// globalApp.Run(":8080")

package illusion

import (
	"errors"
	"regexp"
	//"path/filepath"
)

//存储一个完整调用链
type HandlerInfo struct {
	HttpMethod   string
	RelativePath string
	HandlerChain HandlerChain
}

//切片
type HandlerInfoChain []HandlerInfo

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// Get returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) Get(name string) (string, bool) {
	for _, entry := range ps {
		if entry.Key == name {
			return entry.Value, true
		}
	}
	return "", false
}

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string) (va string) {
	va, _ = ps.Get(name)
	return
}

//每个httpMethod下面会有很多path->HandlerFunc的映射
type methodMap map[string]HandlerFunc

//有什么很好的存储格式???
//这个存储格式太丑了 httpMethod -> MethodMap
type methodHashMap map[string]methodMap

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
	Error error

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

func BluePrint(path, name string) *Blueprint {
	return &Blueprint{
		BasePath:      path,
		Name:          name,
		BeforeChain:   make(HandlerChain, 0, 5), //最大的beforeChain个数
		AfterChain:    make(HandlerChain, 0, 5), //最大的AfterChain个数
		MethodHashMap: make(methodHashMap),      //这个反倒好点
		Error:         nil,
		//illusion:    globalIllusion(),
	}
}

//结合beforeChain + Handler + afterChain形成一个调用链
// methodHashMap形成整个链
func (it *Blueprint) fullChain() HandlerInfoChain {
	if it.Error != nil {
		return nil
	}
	chain := make(HandlerInfoChain, 0, 10)    //还得const来调整
	handlerChain := make(HandlerChain, 0, 20) //这个...
	for httpMethod, urlMap := range it.MethodHashMap {
		//urlMap为url -> handler
		for url, handler := range urlMap {
			handlerChain = append(handlerChain, it.BeforeChain...)
			handlerChain = append(handlerChain, handler)
			handlerChain = append(handlerChain, it.AfterChain...)
			chain = append(chain, HandlerInfo{HttpMethod: httpMethod, RelativePath: url, HandlerChain: handlerChain})
		}
		handlerChain = handlerChain[0:0] //重设一下
	}
	return chain
}

//结合blueprint.BasePath + relativePath形成绝对Url
func (it *Blueprint) truePath(relativePath string) string {
	if it.Error != nil {
		return ""
	}
	return CleanPath(it.BasePath + "/" + relativePath)
}

//我很怀疑append这个操作..
func (it *Blueprint) Before(handler HandlerFunc) *Blueprint {
	if it.Error != nil {
		return nil
	}
	it.BeforeChain = append(it.BeforeChain, handler)
	return it
}

//同上
func (it *Blueprint) After(handler HandlerFunc) *Blueprint {
	if it.Error != nil {
		return nil
	}
	it.AfterChain = append(it.AfterChain, handler)
	return it
}

func (it *Blueprint) handle(httpMethod, relativePath string, handler HandlerFunc) *Blueprint {
	finalPath := it.truePath(relativePath)
	if finalPath == "" {
		it.Error = errors.New("获取绝对路径时出错")
		return it
	}
	//暂时存储在map中
	methodTree := it.MethodHashMap[httpMethod]
	if methodTree == nil {
		methodTree = make(methodMap)
		it.MethodHashMap[httpMethod] = methodTree
	}
	methodTree[finalPath] = handler
	//it.MethodHashMap[httpMethod][finalPath] = handler
	return it
}

func (it *Blueprint) Handle(httpMethod, relativePath string, handler HandlerFunc) *Blueprint {
	if it.Error != nil {
		return nil
	}
	if matches, err := regexp.MatchString("^[A-Z]+$", httpMethod); !matches || err != nil {
		it.Error = err
		panic("http method " + httpMethod + " is not valid")
	}
	return it.handle(httpMethod, relativePath, handler)
}

func (it *Blueprint) Post(relativePath string, handler HandlerFunc) *Blueprint {
	return it.handle("POST", relativePath, handler)
}

func (it *Blueprint) Get(relativePath string, handler HandlerFunc) *Blueprint {
	return it.handle("GET", relativePath, handler)
}

func (it *Blueprint) Delete(relativePath string, handler HandlerFunc) *Blueprint {
	return it.handle("DELETE", relativePath, handler)
}

func (it *Blueprint) Patch(relativePath string, handler HandlerFunc) *Blueprint {
	return it.handle("PATCH", relativePath, handler)
}

func (it *Blueprint) Put(relativePath string, handler HandlerFunc) *Blueprint {
	return it.handle("PUT", relativePath, handler)
}

func (it *Blueprint) Options(relativePath string, handler HandlerFunc) *Blueprint {
	return it.handle("OPTIONS", relativePath, handler)
}

func (it *Blueprint) Head(relativePath string, handler HandlerFunc) *Blueprint {
	return it.handle("HEAD", relativePath, handler)
}

//...
func (it *Blueprint) Any(relativePath string, handler HandlerFunc) *Blueprint {
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

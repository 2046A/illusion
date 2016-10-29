package illusion

import (
	"fmt"
	"net/http"
	"sync"
)

//版本号,什么鬼:)
const Version = "v0.0.1"

//基本处理handler定义
type HandlerFunc func(*Context)
type HandlerChain []HandlerFunc
type MethodTree map[string]*node
type BluePrintTree []*Blueprint

//Illusion应该实现的接口
type IIllusion interface {
	//所有注册在blueprint上的handler会在最终执行前被插入到路由当中
	lazyRegisterAll() *Illusion

	//注册blueprint
	Register(Blueprint) *Illusion

	//设置模板引擎
	//最好不要用，用默认的就好
	SetRender(Render) *Illusion

	//文件相关的加载
	StaticFile(string, string) *Illusion
	Static(string, string) *Illusion
	StaticFS(string, http.FileSystem) *Illusion
}

//核心结构，核心路由器
//这是对httpRouter的改版，参考自gin
//那么Illusion的初始化就是返回一个新的BluePrint ??? 很不错的样子
type Illusion struct {
	//路由器
	render Render

	//最终查找数据存放地
	methodTree MethodTree

	//BluePrint存放地，调用register时blueprint会被临时存放在这里
	//在最终运行前，这个应该被垃圾回收
	//怎样主动回收 ?
	bluePrintTree BluePrintTree

	//Context临时存放地点，随时可以被回收的地点
	//不安全，但能用
	pool sync.Pool

	//下面的待解释...

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 307 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 307 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	//在context中已经被使用了
	//这个...
	ForwardedByClientIP bool
}

/*var defaultIllusion *Illusion
var once sync.Once

func globalIllusion() *Illusion {
	once.Do(func() {
		//defaultIllusion = new(Illusion)
		defaultIllusion = &Illusion{
			render: nil,
			methodTree: make(MethodTree),
			bluePrintTree: make(BluePrintTree, 0, 100),//最大BluePrint的个数，默认100个
			RedirectTrailingSlash: true,
			RedirectFixedPath: false,
			HandleMethodNotAllowed: false,
			ForwardedByClientIP: true,
		}
		defaultIllusion.pool.New = func() *Context{
			return defaultIllusion.allocateContext()
		}
	})
	return defaultIllusion
}*/

//返回一个新的Illusion实例
//鉴于难以抉择，还是把illusion和blueprint的实例化分开为好
func App() (illusion *Illusion) {
	illusion = &Illusion{
		render:                 nil,
		methodTree:             make(MethodTree),
		bluePrintTree:          make(BluePrintTree, 0, 100), //最大BluePrint的个数，默认100个
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      false,
		HandleMethodNotAllowed: false, //并没有使用的一个特性
		ForwardedByClientIP:    true,
	}
	illusion.pool.New = func() interface{} {
		return illusion.allocateContext()
	}
	//b =  Blueprint("/", "default")
	//illusion.Register(b)
	return
}

func (it *Illusion) allocateContext() *Context {
	//Context还没有设计...
	return newContext()
}

func (it *Illusion) Register(bluePrint *Blueprint) *Illusion {
	it.bluePrintTree = append(it.bluePrintTree, bluePrint)
	return it
}

func (it *Illusion) lazyRegisterAll() *Illusion {
	//var tree *node
	for _, bluePrint := range it.bluePrintTree {
		handlerInfoChain := bluePrint.fullChain()
		for _, info := range handlerInfoChain {
			it.addRoute(info.HttpMethod, info.RelativePath, info.HandlerChain)
		}
	}
	return it
}

//设置view目录
func (it *Illusion) ViewPath(path string) *Illusion {
	//todo insert code here
	return it
}

//设置js,css等静态文件目录
func (it *Illusion) Resource(path string) *Illusion {
	return it
}

//添加handle到指定的uriPath
func (it *Illusion) addRoute(httpMethod, uriPath string, handlerChain HandlerChain) {
	tree := it.methodTree[httpMethod]
	if tree == nil {
		tree = new(node)
		it.methodTree[httpMethod] = tree
	}
	tree.addRoute(uriPath, handlerChain)
}

//只有在出错的情况下此函数才会返回
func (it *Illusion) Run(address string) (err error) {
	//address := resolveAddress
	fmt.Println("Listening and serving HTTP on ", address)

	//准备好所有的处理handler
	it.lazyRegisterAll()

	//illusion实现了所有相关的函数
	err = http.ListenAndServe(address, it)
	return
}

//如下应该是作为router必须要实现的一些接口了
//据源代码只是,只要实现这一个接口就好了
//肯定是这样的　:) :)
func (it *Illusion) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	context := it.pool.Get().(*Context)

	//清理原内容
	context.reset()
	context.Request = req
	context.Writer = w

	it.handleRequest(context)

	it.pool.Put(context)
}

//内部真正处理路由
func (it *Illusion) handleRequest(context *Context) {
	httpMethod := context.Request.Method
	path := context.Request.URL.Path

	//找到http方法下面挂着的根
	root := it.methodTree[httpMethod]
	handlers, params, tsr := root.getValue(path)
	if handlers != nil {
		context.handlers = handlers
		context.Params = params
		context.Next()
		//WriteHeaderNow是什么意思啊?
		return
	} else if httpMethod != "CONNECT" && path != "/" {
		if tsr && it.RedirectTrailingSlash {
			redirectTrailingSlash(context)
			return
		}
		if it.RedirectFixedPath && redirectFixedPath(context, root, it.RedirectFixedPath) {
			return
		}
	}

	// TODO: unit test
	// TODO: 这个够不着, 手短
	/*if it.HandleMethodNotAllowed {
		for _, tree := range it.trees {
			if tree.method != httpMethod {
				if handlers, _, _ := tree.root.getValue(path, nil); handlers != nil {
					context.handlers = it.allNoMethod
					serveError(context, 405, default405Body)
					return
				}
			}
		}
	}
	context.handlers = it.allNoRoute
	serveError(context, 404, default404Body)*/
}

func redirectTrailingSlash(c *Context) {
	req := c.Request
	path := req.URL.Path
	code := 301 // Permanent redirect, request with GET method
	if req.Method != "GET" {
		code = 307
	}

	if len(path) > 1 && path[len(path)-1] == '/' {
		req.URL.Path = path[:len(path)-1]
	} else {
		req.URL.Path = path + "/"
	}
	fmt.Print("redirecting request %d: %s --> %s", code, path, req.URL.String())
	http.Redirect(c.Writer, req, req.URL.String(), code)
	//这个WriteHeaderNow是干嘛的 ?
	//c.writermem.WriteHeaderNow()
}

func redirectFixedPath(c *Context, root *node, trailingSlash bool) bool {
	req := c.Request
	path := req.URL.Path

	fixedPath, found := root.findCaseInsensitivePath(
		CleanPath(path),
		trailingSlash,
	)
	if found {
		code := 301 // Permanent redirect, request with GET method
		if req.Method != "GET" {
			code = 307
		}
		req.URL.Path = string(fixedPath)
		fmt.Print("redirecting request %d: %s --> %s", code, path, req.URL.String())
		http.Redirect(c.Writer, req, req.URL.String(), code)
		//这个函数到底是干嘛的 ????
		//c.writermem.WriteHeaderNow()
		return true
	}
	return false
}

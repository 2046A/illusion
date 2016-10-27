package illusion

import (
	"sync"
	"net/http"
	"fmt"
)

//版本号,什么鬼:)
const Version = "v0.0.1"

//基本处理handler定义
type HandlerFunc    func(*Context)
type HandlerChain   []HandlerFunc
type MethodTree     map[string]*node
type BluePrintTree  []*Blueprint

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
	render    Render

	//最终查找数据存放地
	methodTree  MethodTree

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
	ForwardedByClientIP    bool
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
func NewApp() (illusion *Illusion,b *Blueprint) {
	illusion = &Illusion{
		render: nil,
		methodTree: make(MethodTree),
		bluePrintTree: make(BluePrintTree, 0, 100),//最大BluePrint的个数，默认100个
		RedirectTrailingSlash: true,
		RedirectFixedPath: false,
		HandleMethodNotAllowed: false,
		ForwardedByClientIP: true,
	}
	illusion.pool.New = func() interface{}{
		return illusion.allocateContext()
	}
	b =  Blueprint("/", "default")
	illusion.Register(b)
	return
}

func (it *Illusion)allocateContext() *Context{
	//Context还没有设计...
	return new(Context)
}

func (it *Illusion)Register(bluePrint Blueprint) *Illusion{
	it.bluePrintTree = append(it.bluePrintTree, bluePrint)
	return it
}

func (it *Illusion)lazyRegisterAll()*Illusion{
	for _,bluePrint := range it.bluePrintTree{
		//http方法和对应的url->handler
		for httpMethod,methodMap := range bluePrint.MethodHashMap{
			//http方法中的一个tree
			tree := it.methodTree[httpMethod]
			for urlPath,HandlerFunc := range methodMap{
				//把url->handler放入对应的tree中
				tree.addRoute(urlPath, HandlerFunc)
			}
		}
	}
	return it
}

//设置view基础路径
func (it *Illusion)BaseViewPath(path string) *Illusion{
	//todo insert code here
	return it
}

//添加handle到指定的uriPath
func (it *Illusion)addRoute(httpMethod,uriPath string, handler HandlerFunc){
	tree := it.methodTree[httpMethod]
	if tree == nil {
		tree = new(node)
		it.methodTree[httpMethod] = tree
	}
	tree.addRoute(uriPath, handler)
}

//只有在出错的情况下此函数才会返回
func (it *Illusion)Run(address string)(err error){
	//address := resolveAddress
	fmt.Println("Listening and serving HTTP on ", address)

	//准备好所有的处理handler
	it.lazyRegisterAll()
	//illusion实现了所有相关的函数
	err = http.ListenAndServe(address, it)
	return
}

//如下应该是作为router必须要实现的一些接口了
func (it *Illusion)ServeHttp(w http.ResponseWriter,req *http.Request){
	context := it.pool.Get().(*Context)

	it.handleRequest(context)

	it.pool.Put(context)
}

//内部真正处理路由
func (it *Illusion)handleRequest(context *Context){
	httpMethod := context.Request.Method
	path := context.Request.URL.Path

	//找出给定httpMethod的根
	for HttpMethod, root := range it.methodTree{
		//找到根了
		if HttpMethod == httpMethod {
			//找到处理handler
			handlers, params, tsr := root.getValue(path)
			if handlers != nil {
				context.handlers = handlers
				context.Params = params
				context.Next()
				context.writermem.Write
			}
		}
	}

	// Find root of the tree for the given HTTP method
	t := it.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method == httpMethod {
			root := t[i].root
			// Find route in tree
			handlers, params, tsr := root.getValue(path, context.Params)
			if handlers != nil {
				context.handlers = handlers
				context.Params = params
				context.Next()
				context.writermem.WriteHeaderNow()
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
			break
		}
	}

	// TODO: unit test
	if it.HandleMethodNotAllowed {
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
	serveError(context, 404, default404Body)
}



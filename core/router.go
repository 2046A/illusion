package core

import "net/http"

//典型的处理器函数
type Handle func(http.ResponseWriter,*http.Request)


//一个路由器必须实现的接口，实现可插拔的路由器
//不论具体的底层如何实现，只要实现如下接口，即可替换默认路由器
type router interface {

	//GET操作
	Get(string,Handle)

	//HEAD操作
	Head(string,Handle)

	//OPTIONS操作
	Options(string,Handle)

	//POST操作
	Post(string,Handle)

	//PUT操作
	Put(string,Handle)

	//PATCH操作
	Patch(string,Handle)

	//DELETE操作
	Delete(string,Handle)

}

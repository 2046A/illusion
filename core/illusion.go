package core

import "httprouter"

//版本号,:)
const Version = "v0.0.1"

var default404Body = []byte("404 page not found")
var default405Body = []byte("405 method not allowed")


//核心结构
//
type Illusion struct {

	//路由器
	router httprouter.Router
}
package core

import "sync"

//版本号,什么鬼:)
const Version = "v0.0.1"

//基本处理handler定义
type HandlerFunc  func(*Context)
type HandlerChain     []HandlerFunc

//核心结构，核心路由器
//这是对httpRouter的改版，参考自gin
//那么Illusion的初始化就是返回一个新的BluePrint ??? 很不错的样子
type Illusion struct {
	//路由器

}

var defaultIllusion *Illusion
var once  sync.Once
func globalIllusion()*Illusion{
	once.Do(func(){
		defaultIllusion = new(Illusion)
	})
	return defaultIllusion
}

//返回一个新的
func (it *Illusion)New()(bluePrint Blueprint){
	bluePrint = &Blueprint("/", "global")
	return
}

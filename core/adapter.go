package core

import "httprouter"

//用于适配httprouter和框架的一些数据结构

type Param struct {
	httprouter.Param
}

type Params struct {
	httprouter.Params
}

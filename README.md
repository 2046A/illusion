#### illusion
micro web framework for Go, for learning Go

#### feature
+ 基于gin, 主要api兼容gin
+ blueprint主意来自flask,用以取代gin的Group结构，更便于不同路由功能的组织
+ 拥有请求前和请求后接口来组织中间件,会更简单明了

#### DONE
+ 模板引擎
+ 静态资源加载

#### TODO
+ 错误处理 
+ 日志处理
+ 测试

#### example
user.go

```
package main

import (
	"illusion"
	//"net/http"
)

//专门用于abort
func middleware(c *illusion.Context) {
	//c.Write(200, "你说什么")
	//c.Abort()
	c.Append("welcome", "welcome to use illusion framework")
}

func nameEcho(c *illusion.Context) {
	name := c.Param("name")
	if welcome,ok := c.Retrieve("welcome");ok {
		c.Write(200, "your name:" + name + "\t and " + welcome.(string))
	} else {
		c.Write(200, "your name and no welcome")
	}

}

func userBluePrint() *illusion.Blueprint {
	user := illusion.BluePrint("/user", "user")
	user.Before(middleware)
	user.Get("profile/:name", nameEcho)
	return user
}

```
main.go
```
package main

import "illusion"

func main(){
	app := illusion.App()
	//可以直接使用
	ping := illusion.BluePrint("/", "ping")
	ping.Get("ping/:name", func(c *illusion.Context) {
		name := c.Param("name")
		c.Write(200, "hello " + name)
	})
	app.Register(ping)
	//也可以分散到不同的文件中使用
	user := userBluePrint()
	app.Register(user)
	//go  func(){
	app.Run(":8080")
	//}()
}
```
编译运行即可,更多例子见example文件夹

### LICENSE
MIT 


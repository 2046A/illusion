#### illusion
micro web framework for Go, for learning Go

#### feature
+ 基于gin, 主要api兼容gin
+ blueprint主意来自flask,用以取代gin的Group结构,这样不同的路由就可以分散到不同的文件中组织
+ 别告诉我gin也可以，我可是看了源码的

#### TODO
+ 很多 
+ 特别多
+ 你懂的

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
编译运行即可

### LICENSE
MIT 


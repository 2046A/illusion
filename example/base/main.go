package main

import "illusion"

func main(){
	app := illusion.App()
	//可以直接使用
	ping := illusion.BluePrint("/", "ping")
	ping.Get("ping/:name", func(c *illusion.Context) {
		name := c.Param("name")
		c.Echo("hello " + name)
	})
	app.Register(ping)
	//也可以分散到不同的文件中使用
	user := userBluePrint()
	app.Register(user)
	//go  func(){
	app.Run(":8080")
	//}()
}

package main

import (
	//"bytes"
	"illusion"
)
type Title struct {
	Title string
}

func main(){
	app := illusion.App()
	app.ViewPath("example/template/view")
	app.Resource("static")
	index := illusion.BluePrint("/", "index")
	index.Get("/index", func(c *illusion.Context){
		c.View("index.html", Title{Title: "我就是这么吊"})
	})
	ping := illusion.BluePrint("/ping", "ping")
	ping.Get("/", func(c *illusion.Context){
		c.Echo("pong")
	})
	app.Register(index)
	app.Register(ping)
	app.Run(":8080")
}

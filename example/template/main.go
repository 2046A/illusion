package main

import (
	//"bytes"
	"illusion"
//	"net/http"
//	"net/http"
)

type Title struct {
	Title string
}

func notFound(c *illusion.Context){
	c.View("404.html", illusion.TemplateContext{})
}

func main() {
	app := illusion.App()
	app.SetNotFoundHandler(notFound)
	app.ViewPath("example/template/view")
	app.Resource("example/static")
	app.LogPath("example/template/log")
	index := illusion.BluePrint("/", "index")
	index.Before(func(c *illusion.Context){
		if c.HasCookie() {
			sessionId,_ := c.GetSessionId()
			session, err := illusion.Session().StartSession(sessionId)
			if err != nil {
				c.String(200, "不应该出错啊")
				c.Abort()
			}
			session.Store("name", "dean")
		}
	})
	index.Get("/index", func(c *illusion.Context) {
		sessionId, _ := c.GetSessionId()
		session,_ := illusion.Session().StartSession(sessionId)
		val,_ := session.Read("name")
		//c.Status(200)
		//sess,_ := c.Retrieve("sess")
		//val,_ := sess.(illusion.MemorySession).Read("name")
		c.View("index.html", illusion.TemplateContext{"Title": val})
	})
	index.Get("/redirect", func(c *illusion.Context) {
		c.Redirect("/index")
	})
	ping := illusion.BluePrint("/ping", "ping")
	ping.Before(func(c *illusion.Context){
		//illusion.StartSession(c)
		//sess := c.StartSession()
	})
	ping.Get("/", func(c *illusion.Context) {
		c.String(200, "pong too....")
	})
	ping.Get("/string", func(c *illusion.Context) {
		//c.Status(200)
		c.String(200, "pong")
	})
	ping.Get("/json", func(c *illusion.Context) {
		c.Json(200, Title{Title: "I am ok here"})
	})
	ping.Get("/kamila", func(c *illusion.Context) {
		c.Json(200, Title{Title: "Shenm"})
	})
	app.Register(index)
	app.Register(ping)
	app.Run(":8080")
}

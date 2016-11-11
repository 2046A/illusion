package main

import (
	"illusion"
	//"net/http"
)

//专门用于abort
func middleware(c *illusion.Context) {
	c.String(200, "你说什么")
	//c.Abort()
	//c.Append("welcome", "welcome to use illusion framework")
}

func nameEcho(c *illusion.Context) {
	name := c.Param("name")
	if welcome, ok := c.Retrieve("welcome"); ok {
		c.String(200, "your name:" + name + "\t and " + welcome.(string))
	} else {
		c.String(200, "your name and no welcome")
	}

}

func userBluePrint() *illusion.Blueprint {
	user := illusion.BluePrint("/user", "user")
	user.Before(middleware)
	user.Get("profile/:name", nameEcho)
	return user
}

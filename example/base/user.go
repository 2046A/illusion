package main

import (
	"illusion"
	//"net/http"
)

//专门用于abort
func middleware(c *illusion.Context) {
	c.String("你说什么")
	//c.Abort()
	//c.Append("welcome", "welcome to use illusion framework")
}

func nameEcho(c *illusion.Context) {
	name := c.Param("name")
	if welcome, ok := c.Retrieve("welcome"); ok {
		c.String("your name:" + name + "\t and " + welcome.(string))
	} else {
		c.String("your name and no welcome")
	}

}

func userBluePrint() *illusion.Blueprint {
	user := illusion.BluePrint("/user", "user")
	user.Before(middleware)
	user.Get("profile/:name", nameEcho)
	return user
}

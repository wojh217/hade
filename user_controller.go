package main

import (
	"github.com/wojh217/hade/framework/gin"
)

func UserLoginController(c *gin.Context) {
	foo, _ := c.DefaultQueryString("foo", "def")
	c.ISetOkStatus().IJson("ok, UserLoginController: " + foo)
}


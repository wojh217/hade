package main

import (
	"fmt"
	"github.com/wojh217/hade/framework/gin"
	"github.com/wojh217/hade/provider/demo"
)

// subject controller
func SubjectAddController(c *gin.Context) {
	c.ISetOkStatus().IJson( "ok, SubjectAddController")
	return
}

func SubjectListController(c *gin.Context) {
	// 获取实例
	demoService := c.MustMake(demo.Key).(demo.Service)

	// 获取实例提供的方法
	foo := demoService.GetFoo()

	c.ISetOkStatus().IJson(foo)
	return
}

func SubjectDelController(c *gin.Context) {
	c.ISetOkStatus().IJson( "ok, SubjectDelController")
	return
}

func SubjectUpdateController(c *gin.Context) {
	c.ISetOkStatus().IJson( "ok, SubjectUpdateController")
	return
}

func SubjectGetController(c *gin.Context) {
	c.ISetOkStatus().IJson( "ok, SubjectGetController")
	return
}

func SubjectNameController(c *gin.Context) {
	fmt.Println("ok, SubjectNameController")
	c.ISetOkStatus().IJson( "ok, SubjectNameController")
	return
}

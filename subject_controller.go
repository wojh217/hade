package main

import (
	"fmt"
	"github.com/wojh217/hade/framework/gin"
)

// subject controller
func SubjectAddController(c *gin.Context) {
	c.ISetOkStatus().IJson( "ok, SubjectAddController")
	return
}

func SubjectListController(c *gin.Context) {
	c.ISetOkStatus().IJson( "ok, SubjectListController")
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

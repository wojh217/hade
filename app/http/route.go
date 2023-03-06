package http

import (
	"github.com/wojh217/hade/app/http/module/demo"
	"github.com/wojh217/hade/framework/gin"
)

func Routes(r *gin.Engine) {

	// 什么作用？
	r.Static("/dist/", "./dist/")

	demo.Register(r)
}

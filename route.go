package main

import (
	"github.com/wojh217/hade/framework/gin"
)

func registerRouter(core *gin.Engine) {

	core.Use(FuncHandler1())


	//core.Get("/foo", FooControllerHandler)
	core.GET("/login", UserLoginController)

	// 单个url指定中间件, 中间件就是handler格式
	//core.Get("/bar", BarController)
	//core.Get("/panic", PanicController)

	subjectApi := core.Group("/subject")
	{
		subjectApi.DELETE("/:id", SubjectDelController)
		subjectApi.PUT("/:id", SubjectUpdateController)
		subjectApi.GET("/:id", SubjectGetController)
		subjectApi.GET("/list/all", SubjectListController)

		subjectInnerApi := subjectApi.Group("/info")
		{
			subjectInnerApi.GET("/name", SubjectNameController)
		}
	}

	//core.DisplayTree()
}

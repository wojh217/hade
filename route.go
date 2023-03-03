package main

import (
	"fmt"
	"webframework/framework"
)

func registerRouter(core *framework.Core) {
	//core.Use(TimeoutHandler(2 * time.Second))

	//core.Use(RecoverHandler(), CalRequest())
	//core.Use(RecoverHandler())
	//
	//core.Use(CalRequest())
	core.Use(FuncHandler1(), FuncHandler2(), FuncHandler3())
	core.Use(FuncHandler4(), FuncHandler5())

	fmt.Println(core.GetMiddleWares())

	//core.Get("/foo", FooControllerHandler)
	core.Get("/login", UserLoginController)

	// 单个url指定中间件, 中间件就是handler格式
	//core.Get("/bar", BarController)
	//core.Get("/panic", PanicController)

	subjectApi := core.Group("/subject")
	{
		subjectApi.Delete("/:id", SubjectDelController)
		subjectApi.Put("/:id", SubjectUpdateController)
		subjectApi.Get("/:id", SubjectGetController)
		subjectApi.Get("/list/all", SubjectListController)

		subjectInnerApi := subjectApi.Group("/info")
		{
			subjectInnerApi.Get("/name", SubjectNameController)
		}
	}

	//core.DisplayTree()


}

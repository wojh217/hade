package main

import (
	"github.com/wojh217/hade/app/console"
	hadhttp "github.com/wojh217/hade/app/http"
	"github.com/wojh217/hade/framework"
	"github.com/wojh217/hade/framework/provider/app"
	"github.com/wojh217/hade/framework/provider/kernel"
)

func main() {
	// container位于最外层，暴露出来
	// gin的engine、cmd的rootCmd中保存有container
	container := framework.NewHadeContainer()

	// HadeAppProvider属于是框架提供的provider，
	// 其实例本身有Version()、XXXDIr()方法
	container.Bind(&app.HadeAppProvider{BaseFolder: "/tmp"})


	if engine, err := hadhttp.NewHttpEngine(container); err == nil {
		// HadeKernelProvider属于应用提供的provider，
		// 其实例具有HttpEngine()方法
		container.Bind(&kernel.HadeKernelProvider{HttpEngine: engine})
	}

	// 应用程序提供的命令行
	console.RunCommand(container)

}

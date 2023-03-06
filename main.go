package main

import (
	"context"
	hadhttp "github.com/wojh217/hade/app/http"
	"github.com/wojh217/hade/app/provider/demo"
	"github.com/wojh217/hade/framework/gin"
	"github.com/wojh217/hade/framework/middleware"
	"github.com/wojh217/hade/framework/provider/app"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	core := gin.New()

	// 注册两个provider
	core.Bind(&app.HadeAppProvider{BaseFolder: "/tmp"})
	core.Bind(&demo.DemoProvider{})

	// 使用middlerware
	core.Use(gin.Recovery())
	core.Use(middleware.Cost())

	hadhttp.Routes(core)

	server := &http.Server{
		Handler: core,
		Addr:    ":8080",
	}

	go func() {
		server.ListenAndServe()
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal("server shutdown: ", err)
	}

}

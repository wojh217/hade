package main

import (
	"context"
	"github.com/wojh217/hade/framework/gin"
	"github.com/wojh217/hade/provider/demo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)


func main() {
	engine := gin.Default()
	// 注册服务提供者
	engine.Bind(&demo.DemoServiceProvider{})

	registerRouter(engine)

	server := &http.Server{
		Handler: engine,
		Addr:    ":8080",
	}

	go func() {
		server.ListenAndServe()
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<- quit

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal("server shutdown: ", err)
	}

}



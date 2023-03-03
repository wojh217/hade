package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"webframework/framework"
)


func main() {
	core := framework.NewCore()
	registerRouter(core)

	server := &http.Server{
		Handler: core,
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

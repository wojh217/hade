package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"webframework/framework"
)

func Foo2(ctx *framework.Context) error {
	obj := map[string]interface{}{
		"data": nil,
	}
	fooInt := ctx.FormInt("foo", 10)
	obj["data"] = fooInt
	return ctx.Json(http.StatusOK, obj)

}

func FooControllerHandler(ctx *framework.Context) error {
	finish := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)

	durationCtx, cancel := context.WithTimeout(ctx.BaseContext(), 1*time.Second)
	defer cancel()

	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()


		time.Sleep(10 * time.Second)
		ctx.Json(200, "ok")
		finish <- struct{}{}
	}()

	select {
	case <-finish:
		fmt.Println("finish")
	case p := <-panicChan:
		ctx.WriterMux().Lock()
		defer ctx.WriterMux().Unlock()
		log.Println(p)
		ctx.Json(500, "panic")
	case <-durationCtx.Done():
		ctx.WriterMux().Lock()
		defer ctx.WriterMux().Unlock()
		ctx.Json(500, "time out")
		ctx.SetHasTimeout()
	}

	return nil
}
func main() {

	core := framework.NewCore()
	core.Get("/foo", FooControllerHandler)
	server := http.Server{
		Handler: core,
		Addr:    ":8080",
	}
	server.ListenAndServe()

}

package main

import (
	"context"
	"fmt"
	"github.com/wojh217/hade/framework/gin"
	"time"
)



//func FooControllerHandler(ctx *framework.Context) error {
//	finish := make(chan struct{}, 1)
//	panicChan := make(chan interface{}, 1)
//
//	durationCtx, cancel := context.WithTimeout(ctx.BaseContext(), 1*time.Second)
//	defer cancel()
//
//	go func() {
//		defer func() {
//			if p := recover(); p != nil {
//				panicChan <- p
//			}
//		}()
//
//
//		time.Sleep(10 * time.Second)
//		ctx.SetStatus(200).Json( "ok")
//		finish <- struct{}{}
//	}()
//
//	select {
//	case <-finish:
//		fmt.Println("finish")
//	case p := <-panicChan:
//		ctx.WriterMux().Lock()
//		defer ctx.WriterMux().Unlock()
//		log.Println(p)
//		ctx.SetStatus(500).Json( "panic")
//	case <-durationCtx.Done():
//		ctx.WriterMux().Lock()
//		defer ctx.WriterMux().Unlock()
//		ctx.SetStatus(500).Json( "time out")
//		ctx.SetHasTimeout()
//	}
//
//	return nil
//}


func PanicController(c *gin.Context) error {
	ch := make(chan int)
	close(ch)
	ch <- 1

	c.ISetStatus(200).IJson( "ok, PanicController")
	return nil
}



func BarController(c *gin.Context) error {
	c.ISetStatus(200).IJson( "ok, BarController")
	return nil
}



// 自定义中间件
// 超时中间件， 接收一个时间参数，返回handler
func TimeoutHandler(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context)  {
		fmt.Printf("time.Duration: %d\n", d)

		durationCtx, cancel := context.WithTimeout(c, d)
		defer cancel()

		c.Request.WithContext(durationCtx)
		finish := make(chan bool)

		go func() {
			// 真正的操作
			c.Next()

			finish <- true
		}()

		select {
		case <- finish:
		case <- durationCtx.Done():
			c.ISetStatus(500).IJson( "time out")
		}
	}
}


// RecoverHandler
func RecoverHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if p := recover(); p != nil {
				fmt.Printf("catch panic: %v\n", p)
				//c.SetStatus(500).Json( fmt.Sprintf("catch panic: %v", p)) // 这里直接回复报错
			}
		}()

		// 调用的途中遇到了panic
		c.Next()
	}
}

func FuncHandler1() gin.HandlerFunc {
	return func(c *gin.Context)  {
		fmt.Println("FuncHandler1")

		// 没有返回值
		c.Next()
	}
}



// 统计请求的中间件
func CalRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		begintime := time.Now()
		defer func() {
			fmt.Printf("request uri: %s, begintime: %s, endtime: %s, cost: %v\n", c.Request.URL, begintime, time.Now(), time.Now().Sub(begintime))
		}()
		c.Next()
	}
}


package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"webframework/framework"
)



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
		ctx.SetStatus(200).Json( "ok")
		finish <- struct{}{}
	}()

	select {
	case <-finish:
		fmt.Println("finish")
	case p := <-panicChan:
		ctx.WriterMux().Lock()
		defer ctx.WriterMux().Unlock()
		log.Println(p)
		ctx.SetStatus(500).Json( "panic")
	case <-durationCtx.Done():
		ctx.WriterMux().Lock()
		defer ctx.WriterMux().Unlock()
		ctx.SetStatus(500).Json( "time out")
		ctx.SetHasTimeout()
	}

	return nil
}


func PanicController(c *framework.Context) error {
	ch := make(chan int)
	close(ch)
	ch <- 1

	c.SetStatus(200).Json( "ok, PanicController")
	return nil
}



func BarController(c *framework.Context) error {
	c.SetStatus(200).Json( "ok, BarController")
	return nil
}



// 自定义中间件
// 超时中间件， 接收一个时间参数，返回handler
func TimeoutHandler(d time.Duration) framework.ControllerHandler {
	return func(c *framework.Context) error {
		fmt.Printf("time.Duration: %d\n", d)

		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		c.GetRequest().WithContext(durationCtx)
		finish := make(chan error)

		go func() {
			finish <- c.Next()
		}()

		select {
		case err := <- finish:
			return err
		case <- durationCtx.Done():
			c.WriterMux().Lock()
			defer c.WriterMux().Unlock()
			c.SetStatus(500).Json( "time out")
			c.SetHasTimeout()
			return nil
		}
	}
}



func RecoverHandler() framework.ControllerHandler {
	return func(c *framework.Context) error {
		defer func() {
			if p := recover(); p != nil {
				fmt.Printf("catch panic: %v\n", p)
				c.SetStatus(500).Json( fmt.Sprintf("catch panic: %v", p)) // 这里直接回复报错
			}
		}()

		// 调用的途中遇到了panic
		return c.Next()
	}
}

func FuncHandler1() framework.ControllerHandler {
	return func(c *framework.Context) error {
		fmt.Println("FuncHandler1")

		return c.Next()
	}
}

func FuncHandler2() framework.ControllerHandler {
	return func(c *framework.Context) error {
		fmt.Println("FuncHandler2")

		//ch := make(chan bool)
		//close(ch)
		//
		//ch <- true
		return c.Next()
	}
}

func FuncHandler3() framework.ControllerHandler {
	return func(c *framework.Context) error {
		fmt.Println("FuncHandler3")

		return c.Next()
	}
}
func FuncHandler4() framework.ControllerHandler {
	return func(c *framework.Context) error {
		fmt.Println("FuncHandler4")

		return c.Next()
	}
}
func FuncHandler5() framework.ControllerHandler {
	return func(c *framework.Context) error {
		fmt.Println("FuncHandler5")

		return c.Next()
	}
}

// 统计请求的中间件
func CalRequest() framework.ControllerHandler {
	return func(c *framework.Context) error {
		begintime := time.Now()
		defer func() {
			fmt.Printf("request uri: %s, begintime: %s, endtime: %s, cost: %v\n", c.GetRequest().URL, begintime, time.Now(), time.Now().Sub(begintime))
		}()

		return c.Next()
	}
}


package command

import (
	"context"
	"fmt"
	"github.com/wojh217/hade/framework/cobra"
	"github.com/wojh217/hade/framework/contract"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// initAppCommand 初始化app命令和其子命令
func initAppCommand() *cobra.Command {
	appCommand.AddCommand(appStartCommand)
	return appCommand
}

// AppCommand 是命令行参数第一级为 app 的命令，它没有实际功能，只是打印帮助文档
var appCommand = &cobra.Command{
	Use:   "app",
	Short: "业务应用控制命令",
	Long:  "业务应用控制命令，其包含业务启动，关闭，重启，查询等功能",
	RunE: func(c *cobra.Command, args []string) error {
		// 打印帮助文档
		c.Help()
		return nil
	},
}

// appStartCommand 启动一个Web服务
var appStartCommand = &cobra.Command{
	Use: "start",
	Short: "启动一个web服务",
	RunE: func(c *cobra.Command, args []string) error {
		// 使用start子命令才真正启动http服务
		// 要做到这一点，需要把engine写到一个service注册到container，并把此container保存到rootCmd
		container := c.GetContainer()
		fmt.Printf("start get container: %v\n", container)
		kernelService := container.MustMake(contract.KernelKey).(contract.Kernel)
		core := kernelService.HttpEngine()

		// 创建一个Server服务
		server := &http.Server{
			Handler: core,
			Addr:    ":8888",
		}

		// 这个goroutine是启动服务的goroutine
		go func() {
			server.ListenAndServe()
		}()

		// 当前的goroutine等待信号量
		quit := make(chan os.Signal)
		// 监控信号：SIGINT, SIGTERM, SIGQUIT
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		// 这里会阻塞当前goroutine等待信号
		<-quit

		// 调用Server.Shutdown graceful结束
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(timeoutCtx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}

		return nil
	},
}
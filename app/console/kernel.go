package console

import (
	"github.com/wojh217/hade/app/console/command/demo"
	"github.com/wojh217/hade/framework"
	"github.com/wojh217/hade/framework/cobra"
	"github.com/wojh217/hade/framework/command"
)

//RunCommand is command
// rootCmd只打印help
func RunCommand(container framework.Container) error {
	var rootCmd = &cobra.Command{
		Use:   "hade",
		Short: "hade 命令",
		Long:  "hade 框架提供的命令行工具",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 执行help函数，什么操作都没有，真正运行程序放在start子命令中
			cmd.InitDefaultHelpFlag()
			return cmd.Help()
		},
		// 不需要出现cobra默认的completion子命令
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	// 把container写到rootCmd
	rootCmd.SetContainer(container)

	// 使用框架中添加子命令 demo  和 app start
	command.AddKernelCommands(rootCmd)

	// 使用应用程序添加子命令foo
	AddAppCommand(rootCmd)

	return rootCmd.Execute()
}

func AddAppCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(demo.InitFoo())
}

package demo

import (
	"github.com/wojh217/hade/framework/cobra"
	"log"
)

func InitFoo() *cobra.Command {
	FooCommand.AddCommand(Foo1Command)
	return FooCommand
}

var FooCommand = &cobra.Command{
	Use: "foo",
	Short: "foo的简要说明",
	Long: "foo的长说明",
	Aliases: []string{"fo", "f"},
	Example: "foo命令的例子",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		log.Printf("foo命令捕获到的container: %v\n", container)
		return nil
	},
}

var Foo1Command = &cobra.Command{
	Use: "foo1",
	Short: "foo1的简短说明",
	Long: "foo1的长说明",
	Aliases: []string{"fo1", "f1"},
	Example: "foo1命令的例子",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		log.Printf("foo1命令捕获到的container: %v\n", container)
		return nil
	},

}

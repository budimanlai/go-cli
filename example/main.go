package main

import (
	"fmt"

	gocli "github.com/budimanlai/go-cli"
	"github.com/budimanlai/go-cli/example/commands"
)

func main() {
	cli := gocli.NewCliWithConfig(gocli.CliOptions{
		AppName:     "Go CLI Demo",
		Version:     "1.0.0",
		ConfigFile:  []string{"config/main.conf"},
		RuntimePath: "runtime/",
	})

	cli.AddCommand("random_string", commands.RandomString)
	cli.AddCommand("world", commands.World)
	cli.AddCommand("clean_log", commands.Clearlog)
	cli.AddCommand("listen", commands.Listen)
	cli.AddCommand("run", commands.Listen)

	e := cli.Run()
	if e != nil {
		fmt.Println(e.Error())
	}
}

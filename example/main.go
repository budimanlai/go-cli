package main

import (
	"fmt"

	gocli "github.com/budimanlai/go-cli"
	"github.com/budimanlai/go-cli/example/commands"
)

func main() {
	cli := gocli.NewCliWithConfig(gocli.CliOptions{
		ConfigFile:  []string{"config/main.conf"},
		RuntimePath: "output/",
	})
	cli.AppName = "Go CLI Demo"
	cli.Version = "1.0.0"

	cli.AddCommand("random_string", commands.RandomString)
	cli.AddCommand("world", commands.World)
	cli.AddCommand("clean_log", commands.Clearlog)
	cli.AddCommand("listen", commands.Listen)

	e := cli.Run()
	if e != nil {
		fmt.Println(e.Error())
	}
}

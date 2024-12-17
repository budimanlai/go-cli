package main

import (
	"fmt"

	gocli "github.com/budimanlai/go-cli"
)

func main() {
	cli := gocli.NewCli()
	cli.AppName = "Go CLI Demo"
	cli.Version = "1.0.0"

	cli.AddCommand("helo", func(c *gocli.Cli) {
		fmt.Println("Helo CLI Demo")
	})

	cli.AddCommand("world", func(c *gocli.Cli) {
		fmt.Println("World CLI Demo")
		fmt.Println(c.Args.GetRawArgs())
		fmt.Println("Port:", c.Args.GetInt("port"))
	})

	e := cli.Run()
	if e != nil {
		fmt.Println(e.Error())
	}
}

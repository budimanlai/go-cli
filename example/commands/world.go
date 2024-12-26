package commands

import (
	"fmt"

	gocli "github.com/budimanlai/go-cli"
)

func World(c *gocli.Cli) {
	c.LoadConfig()

	fmt.Println("World CLI Demo")
	fmt.Println(c.Args.GetRawArgs())
	fmt.Println("Port:", c.Args.GetInt("port"))
	fmt.Println("database:", c.Config.GetString("database.database"))
	fmt.Println("hostname:", c.Config.GetString("database.hostname"))
	fmt.Println("username:", c.Config.GetString("database.username"))
	fmt.Println("port:", c.Config.GetInt("database.port"))
	fmt.Println("runtime path:", c.RuntimePath())
}

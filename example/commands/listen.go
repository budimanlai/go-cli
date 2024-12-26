package commands

import (
	gocli "github.com/budimanlai/go-cli"
)

func Listen(c *gocli.Cli) {
	c.Log("Listen command example")
	c.Log("Press ctrl+c to stop the service")

	number := c.Args.GetIntOr("number", 1)

	c.Listen(gocli.CliListen{
		TimeLoop: 4,
		OnLoop: func() {
			c.Log("Service in loop... ", number)
			number = number + 1
		},
		// OnShutdown: func() {
		// 	c.Log("Shutdown service detected...")
		// 	c.Log("Cache cleanup....")
		// 	time.Sleep(3 * time.Second)

		// 	c.Log("Disconnect database...")
		// 	time.Sleep(1 * time.Second)

		// 	c.Log("Other process before stop the service...")
		// 	time.Sleep(4 * time.Second)

		// 	c.Log("Ok, now stop the service...")
		// },
	})
}

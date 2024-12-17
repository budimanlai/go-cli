package gocli

import (
	"fmt"

	goargs "github.com/budimanlai/go-args"
)

type CliHandler = func(c *Cli)

type Cli struct {
	AppName string
	Version string

	Args    *goargs.Args
	Config  Config
	handler map[string]CliHandler
}

func NewCli() *Cli {
	c := &Cli{}

	c.Args = &goargs.Args{}
	c.Args.Parse()

	c.handler = map[string]CliHandler{}
	c.addDefaultCommand()

	return c
}

func (c *Cli) addDefaultCommand() {
	c.AddCommandAndAlias("version", "v", func(c *Cli) {
		fmt.Println(c.AppName, "\nVersi", c.Version)
	})
}

func (c *Cli) AddCommandAndAlias(command string, alias string, handler CliHandler) {
	c.handler[command] = handler
	c.handler[alias] = handler
}

func (c *Cli) AddCommand(command string, handler CliHandler) {
	c.handler[command] = handler
}

func (c *Cli) Run() error {
	if h, e := c.handler[c.Args.Command]; e {
		h(c)
	} else {
		fmt.Print(fmt.Sprintf("Command '%s' not found", c.Args.Command) + "\n")
	}

	return nil
}

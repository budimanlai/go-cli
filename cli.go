package gocli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	goargs "github.com/budimanlai/go-args"
	goconfig "github.com/budimanlai/go-config"
)

type CliHandler = func(c *Cli)

type Cli struct {
	AppName    string
	Version    string
	IsShutdown bool

	Args        *goargs.Args
	Config      *goconfig.Config
	handler     map[string]CliHandler
	runtimePath string
	configFile  []string
}

type CliListen struct {
	OnLoop     func()
	OnShutdown func()
	TimeLoop   time.Duration
}

const (
	YYYYMMDDHHMMSS string = "2006-01-02 15:04:05"
)

func NewCli() *Cli {
	c := NewCliWithConfig(CliOptions{})

	return c
}

func NewCliWithConfig(config CliOptions) *Cli {
	c := &Cli{}

	c.Config = &goconfig.Config{}
	c.configFile = config.ConfigFile

	if config.AutoLoadConfig {
		if len(c.configFile) != 0 {
			e := c.Config.Open(config.ConfigFile...)

			if e != nil {
				panic(e)
			}
		}
	}

	c.Args = &goargs.Args{}
	c.Args.Parse()

	if config.RuntimePath == "" {
		config.RuntimePath = "runtime/"
	}
	c.runtimePath = config.RuntimePath

	c.handler = map[string]CliHandler{}
	c.addDefaultCommand()

	return c
}

func (c *Cli) LoadConfig() {
	if len(c.configFile) != 0 {
		e := c.Config.Open(c.configFile...)

		if e != nil {
			panic(e)
		}
	}
}

func (c *Cli) RuntimePath() string { return c.runtimePath }

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

func (c *Cli) Log(a ...interface{}) {
	now := time.Now()
	date := now.Format(YYYYMMDDHHMMSS)
	fmt.Print("[" + date + "] ")
	fmt.Println(a...)
}

func (c *Cli) Listen(handler CliListen) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {
		sig := <-sigs

		signal.Stop(sigs)
		close(sigs)

		c.IsShutdown = true

		c.Log(sig)
		if handler.OnShutdown != nil {
			handler.OnShutdown()
			c.Log("Service stoped...")
			os.Exit(0)
		}

		done <- true
	}()

	c.IsShutdown = false
	if handler.OnLoop != nil {
		if handler.TimeLoop < 2 {
			handler.TimeLoop = 2
		}
		for {
			if c.IsShutdown {
				<-done
				break
			}

			handler.OnLoop()
			time.Sleep(handler.TimeLoop * time.Second)
		}
	}
}

func (c *Cli) Run() error {
	if h, e := c.handler[c.Args.Command]; e {
		h(c)
	} else {
		fmt.Print(fmt.Sprintf("Command '%s' not found", c.Args.Command) + "\n")
	}

	return nil
}

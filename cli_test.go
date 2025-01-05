package gocli

import (
	"syscall"
	"testing"
	"time"

	goargs "github.com/budimanlai/go-args"
	"github.com/stretchr/testify/assert"
)

func TestNewCli(t *testing.T) {
	cli := NewCli()
	assert.NotNil(t, cli)
	assert.NotNil(t, cli.Config)
	assert.False(t, cli.IsShutdown)
}

func TestLoadConfig(t *testing.T) {
	cli := NewCli()
	cli.configFile = []string{"example/config/main.conf"}

	// Mock the Open function
	// cli.Config = &goconfig.Config{
	// 	Open: func(files ...string) error {
	// 		return nil
	// 	},
	// }

	assert.NotPanics(t, func() {
		cli.LoadConfig()
	})
}

func TestAddCommandAndAlias(t *testing.T) {
	cli := NewCli()
	cli.handler = make(map[string]CliHandler)

	cli.AddCommandAndAlias("test", "t", func(c *Cli) {
		c.Log("Test command executed")
	})

	assert.Contains(t, cli.handler, "test")
	assert.Contains(t, cli.handler, "t")
}

func TestLog(t *testing.T) {
	cli := NewCli()
	assert.NotPanics(t, func() {
		cli.Log("This is a test log")
	})
}

func TestListen(t *testing.T) {
	cli := NewCli()
	cli.handler = make(map[string]CliHandler)

	listenHandler := CliRunLoop{
		OnLoop: func() {
			cli.Log("Looping...")
		},
		OnShutdown: func() {
			cli.Log("Shutting down...")
		},
		TimeLoop: 1 * time.Second,
	}

	go func() {
		time.Sleep(2 * time.Second)
		cli.IsShutdown = true
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	assert.NotPanics(t, func() {
		cli.RunLoop(listenHandler)
	})
}

func TestRun(t *testing.T) {
	cli := NewCli()
	cli.handler = make(map[string]CliHandler)
	cli.Args = &goargs.Args{Command: "test"}

	cli.AddCommand("test", func(c *Cli) {
		c.Log("Test command executed")
	})

	assert.NotPanics(t, func() {
		err := cli.Run()
		assert.Nil(t, err)
	})
}

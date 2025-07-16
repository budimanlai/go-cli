package gocli

import (
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

	// Track if OnLoop was called
	loopCalled := false
	loopCount := 0

	listenHandler := CliRunLoop{
		OnLoop: func() {
			cli.Log("Looping...")
			loopCalled = true
			loopCount++

			// After a few loops, set shutdown to avoid infinite running
			if loopCount >= 2 {
				cli.IsShutdown = true
			}
		},
		OnShutdown: func() {
			cli.Log("Shutting down...")
		},
		TimeLoop: 50 * time.Millisecond, // Very short interval for fast test
	}

	// Create a done channel to signal completion
	done := make(chan bool, 1)

	// Run the loop in a goroutine and handle panic gracefully
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Expected panic from os.Exit(), signal completion
				t.Logf("Caught expected panic: %v", r)
				done <- true
			}
		}()
		cli.RunLoop(listenHandler)
		done <- true // In case it exits normally (shouldn't happen)
	}()

	// Wait for completion or timeout
	select {
	case <-done:
		// Test completed (either normally or via panic)
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out - RunLoop didn't complete")
	}

	// Verify that the loop was called
	assert.True(t, loopCalled, "OnLoop should have been called at least once")
	assert.GreaterOrEqual(t, loopCount, 1, "Loop should have run at least once")
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

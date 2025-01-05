package gocli

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"

	goargs "github.com/budimanlai/go-args"
	goconfig "github.com/budimanlai/go-config"
)

type CliHandler = func(c *Cli)

type Cli struct {
	IsShutdown bool

	Args        *goargs.Args
	Config      *goconfig.Config
	handler     map[string]CliHandler
	runtimePath string
	configFile  []string
	appName     string
	version     string
	wg          sync.WaitGroup // WaitGroup for task synchronization
	mutex       sync.Mutex     // Mutex for shutdown status
	pidFile     string         // File to store PID
}

type CliRunLoop struct {
	OnLoop     func()
	OnShutdown func()
	TimeLoop   time.Duration
}

const (
	YYYYMMDDHHMMSS string = "2006-01-02 15:04:05"
)

// NewCli creates a new instance of Cli with default options.
func NewCli() *Cli {
	c := NewCliWithConfig(CliOptions{})

	return c
}

// NewCliWithConfig creates a new instance of Cli with the provided options.
//
// Parameters:
// - config: CliOptions struct containing configuration options for the CLI.
func NewCliWithConfig(config CliOptions) *Cli {
	c := &Cli{}

	c.appName = config.AppName
	c.version = config.Version

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

	// Get the binary file name
	binaryName := filepath.Base(os.Args[0])
	c.pidFile = config.RuntimePath + binaryName + ".pid"

	c.handler = map[string]CliHandler{}
	c.addDefaultCommand()

	return c
}

// LoadConfig loads the configuration files specified in the Cli instance.
func (c *Cli) LoadConfig() {
	if len(c.configFile) != 0 {
		e := c.Config.Open(c.configFile...)

		if e != nil {
			panic(e)
		}
	}
}

// RuntimePath returns the runtime path of the Cli instance.
func (c *Cli) RuntimePath() string { return c.runtimePath }

// addDefaultCommand adds the default commands to the CLI instance.
func (c *Cli) addDefaultCommand() {
	c.AddCommandAndAlias("version", "v", func(c *Cli) {
		fmt.Println(c.appName, "\nVersi", c.version)
	})
}

// AddCommandAndAlias adds a command and its alias to the CLI instance.
//
// Parameters:
// - command: The main command string.
// - alias: The alias for the command.
// - handler: The function to handle the command.
func (c *Cli) AddCommandAndAlias(command string, alias string, handler CliHandler) {
	c.handler[command] = handler
	c.handler[alias] = handler
}

// AddCommand adds a command to the CLI instance.
//
// Parameters:
// - command: The command string.
// - handler: The function to handle the command.
func (c *Cli) AddCommand(command string, handler CliHandler) {
	c.handler[command] = handler
}

// Log logs a message with a timestamp.
//
// Parameters:
// - a: The message to log.
func (c *Cli) Log(a ...interface{}) {
	now := time.Now()
	date := now.Format(YYYYMMDDHHMMSS)
	fmt.Print("[" + date + "] ")
	fmt.Println(a...)
}

// listenSignal listens for system signals and handles them accordingly.
//
// Parameters:
// - handler: A CliRunLoop struct containing the OnLoop and OnShutdown functions and the TimeLoop duration.
func (c *Cli) listenSignal(handler CliRunLoop) {
	// Channel to capture system signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP) // Add SIGHUP for reload

	// Goroutine to handle signals
	go func() {
		defer close(sigs) // Ensure the channel is closed
		for sig := range sigs {
			switch sig {
			case syscall.SIGHUP: // Reload configuration
				c.Log("Configuration reload signal received.")
				go func() {
					c.Log("Reloading configuration...")
					tempConfig := &goconfig.Config{}
					err := tempConfig.Open(c.configFile...)
					if err != nil {
						c.Log("Failed to load configuration:", err)
						return
					}

					newConfig := *tempConfig // Create a copy outside the mutex
					c.mutex.Lock()
					c.Config = &newConfig
					c.mutex.Unlock()
					c.Log("Configuration successfully reloaded.")
				}()

			case syscall.SIGINT, syscall.SIGTERM: // Shutdown
				// Mark shutdown
				c.mutex.Lock()
				c.IsShutdown = true
				c.mutex.Unlock()

				// Execute OnShutdown handler if present
				if handler.OnShutdown != nil {
					c.wg.Add(1)

					// Execute handler in a separate goroutine
					go func() {
						defer c.wg.Done()
						defer func() {
							if r := recover(); r != nil {
								c.Log("Panic detected in OnShutdown:", r)
							}
						}()
						c.Log("Start shutdown...")
						handler.OnShutdown()
						c.Log("Shutdown done...")
					}()
				}

				c.wg.Wait()
				signal.Stop(sigs) // Stop receiving signals

				os.Exit(0) // Exit after all tasks are completed
			}
		}
	}()
}

// RunLoop starts the main loop for the CLI application.
// It listens for system signals and executes the provided handler functions.
//
// Parameters:
// - handler: A CliRunLoop struct containing the OnLoop and OnShutdown functions and the TimeLoop duration.
//
// The loop will execute the OnLoop function at intervals specified by TimeLoop.
// If a SIGINT or SIGTERM signal is received, the OnShutdown function will be executed and the application will exit.
// If a SIGHUP signal is received, the configuration will be reloaded.
func (c *Cli) RunLoop(handler CliRunLoop) {
	c.listenSignal(handler)

	if handler.TimeLoop < 2 {
		handler.TimeLoop = 2 * time.Second // Default minimum 2 seconds
	}
	ticker := time.NewTicker(handler.TimeLoop)
	defer ticker.Stop() // Ensure ticker stops on exit

	// Main loop to run tasks
	c.IsShutdown = false
	if handler.OnLoop != nil {
		for range ticker.C { // Using for range
			c.mutex.Lock()
			if c.IsShutdown {
				c.mutex.Unlock()
				c.Log("Waiting shutdown")
				c.wg.Wait()
				return
			}
			c.mutex.Unlock()

			c.wg.Add(1) // Task synchronization
			go func() {
				defer c.wg.Done() // Ensure Done is called even if OnLoop panics
				handler.OnLoop()  // Process task (blocking)
			}()
		}
	}
}

// Run executes the command specified in the CLI arguments.
//
// Returns an error if the command is not found.
func (c *Cli) Run() error {
	if h, exists := c.handler[c.Args.Command]; exists {
		h(c)
	} else {
		fmt.Printf("Command '%s' not found\n", c.Args.Command)
	}

	return nil
}

// StartService adds a command to start a service and a command to start the service as a daemon.
//
// Parameters:
// - command: The command to start the service.
// - startCmd: The command to start the service as a daemon.
// - handler: The function to handle the service command.
func (c *Cli) StartService(command string, startCmd string, handler CliHandler) {
	c.AddCommand(command, handler)
	c.AddCommand(startCmd, func(c *Cli) {
		c.startDaemon(command)
	})

}

// startDaemon starts the specified command as a daemon process.
//
// Parameters:
// - command: The command to start as a daemon.
func (c *Cli) startDaemon(command string) {
	// Get the binary file name
	binaryName := filepath.Base(os.Args[0])
	logFileName := c.runtimePath + binaryName + ".log"

	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	defer logFile.Close()

	// Include additional arguments
	args := append([]string{command}, c.Args.GetRawArgs()[1:]...)

	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} // Detach the process
	err = cmd.Start()
	if err != nil {
		fmt.Println("Failed to start daemon:", err)
		return
	}
	err = os.WriteFile(c.pidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
	if err != nil {
		fmt.Println("Failed to write PID file:", err)
		return
	}
	os.Exit(0)
}

// StopService adds a command to stop a running daemon service.
//
// Parameters:
// - stopCmd: The command to stop the daemon service.
func (c *Cli) StopService(stopCmd string) {
	c.AddCommand(stopCmd, func(c *Cli) {
		c.stopDaemon()
	})
}

// stopDaemon stops the running daemon process.
func (c *Cli) stopDaemon() {
	pidData, err := os.ReadFile(c.pidFile)
	if err != nil {
		fmt.Println("Failed to read PID file:", err)
		return
	}
	pid, err := strconv.Atoi(string(pidData))
	if err != nil {
		fmt.Println("Invalid PID:", err)
		return
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("Failed to find process:", err)
		return
	}
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		fmt.Println("Failed to stop process:", err)
		return
	}
	os.Remove(c.pidFile)
}

# Go-CLI

[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/budimanlai/go-cli.svg)](https://pkg.go.dev/github.com/budimanlai/go-cli)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

Go-CLI is a powerful and easy-to-use library for creating command-line interface (CLI) applications in Go. This library provides comprehensive features to simplify CLI development, including argument parsing, sub-commands, configuration management, background services, and automatic help generation.

## Features

- 🚀 **Easy CLI Creation**: Simple API to create CLI applications quickly
- 📝 **Argument Parsing**: Built-in support for command-line arguments and options
- 🔧 **Configuration Management**: Load configuration from files with hot-reload support
- 🏃 **Sub-commands**: Support for multiple commands and aliases
- 🔄 **Background Services**: Run services with graceful shutdown
- 📊 **Logging**: Built-in logging with configurable output
- 🛡️ **Signal Handling**: Proper signal handling for graceful shutdown
- 🔧 **Runtime Management**: PID file management and process control

## Installation

```bash
go get github.com/budimanlai/go-cli
```

## Dependencies

- `github.com/budimanlai/go-args v0.0.1` - Command-line argument parsing
- `github.com/budimanlai/go-config v0.0.4` - Configuration management

## Quick Start

### Basic CLI Application

```go
package main

import (
    "fmt"
    gocli "github.com/budimanlai/go-cli"
)

func main() {
    // Create a new CLI instance
    cli := gocli.NewCli()
    
    // Add a simple command
    cli.AddCommand("hello", func(c *gocli.Cli) {
        c.Log("Hello, World!")
    })
    
    // Run the CLI
    if err := cli.Run(); err != nil {
        fmt.Println(err.Error())
    }
}
```

### Advanced CLI with Configuration

```go
package main

import (
    "fmt"
    gocli "github.com/budimanlai/go-cli"
    "github.com/budimanlai/go-cli/example/commands"
)

func main() {
    cli := gocli.NewCliWithConfig(gocli.CliOptions{
        AppName:     "My CLI App",
        Version:     "1.0.0",
        ConfigFile:  []string{"config/main.conf"},
        RuntimePath: "runtime/",
    })

    // Add commands
    cli.AddCommand("random_string", commands.RandomString)
    cli.AddCommand("world", commands.World)
    cli.AddCommand("clean_log", commands.Clearlog)
    
    // Add commands with aliases
    cli.AddCommandAndAlias("status", "s", func(c *gocli.Cli) {
        c.Log("Application is running")
    })
    
    // Start background service
    cli.StartService("run", "mulai", commands.Listen)
    cli.StopService("stop")

    if err := cli.Run(); err != nil {
        fmt.Println(err.Error())
    }
}
```

### Background Service with Loop

```go
func Listen(c *gocli.Cli) {
    handler := gocli.CliRunLoop{
        OnLoop: func() {
            c.Log("Service is running...")
            // Your service logic here
        },
        OnShutdown: func() {
            c.Log("Service shutting down...")
            // Cleanup logic here
        },
        TimeLoop: 5 * time.Second, // Run every 5 seconds
    }
    
    c.RunLoop(handler)
}
```

## API Reference

### Core Types

#### `Cli` struct
The main CLI application instance.

#### `CliOptions` struct
Configuration options for CLI initialization:
```go
type CliOptions struct {
    AppName     string   // Application name
    Version     string   // Application version
    ConfigFile  []string // Configuration file paths
    RuntimePath string   // Runtime directory path
}
```

#### `CliRunLoop` struct
Configuration for background service loops:
```go
type CliRunLoop struct {
    OnLoop     func()        // Function called in each loop iteration
    OnShutdown func()        // Function called on graceful shutdown
    TimeLoop   time.Duration // Interval between loop iterations
}
```

### Main Functions

- `NewCli() *Cli` - Create a new CLI instance with default settings
- `NewCliWithConfig(options CliOptions) *Cli` - Create CLI with custom configuration
- `AddCommand(command string, handler CliHandler)` - Add a command
- `AddCommandAndAlias(command, alias string, handler CliHandler)` - Add command with alias
- `StartService(command, alias string, handler CliHandler)` - Add background service command
- `StopService(command string)` - Add stop service command
- `Run() error` - Run the CLI application
- `RunLoop(handler CliRunLoop)` - Run background service loop
- `LoadConfig()` - Load configuration from files
- `Log(message ...interface{})` - Log messages with timestamp

## Example Commands

The library includes example commands in the `example/commands/` directory:

- **random_string**: Generate random strings
- **world**: Display greeting message  
- **clear_log**: Clear application logs
- **listen**: Background service example

## Configuration

The library supports configuration files through the `go-config` dependency. You can specify configuration files when creating the CLI instance:

```go
cli := gocli.NewCliWithConfig(gocli.CliOptions{
    ConfigFile: []string{"config/app.conf", "config/database.conf"},
})
```

## Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestNewCli
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Changelog

### v0.0.4 (Latest)
- Updated go-config dependency to v0.0.4
- Improved test coverage and reliability
- Enhanced error handling
- Better signal handling for graceful shutdown

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Links

- [GitHub Repository](https://github.com/budimanlai/go-cli)
- [Go Package Documentation](https://pkg.go.dev/github.com/budimanlai/go-cli)
- [Examples](https://github.com/budimanlai/go-cli/tree/main/example)
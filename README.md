# Go-CLI

Go-CLI is a library for easily creating command-line interface (CLI) applications using the Go programming language. This library provides various features to simplify the creation of CLIs, such as argument parsing, sub-commands, and automatic help.

## Usage
This library is useful for:
- Quickly and easily creating CLI applications.
- Managing arguments and options from the command line.
- Providing automatic help for CLI users.

## How to Use
1. Installation:
   ```sh
   go get github.com/budimanlai/go-cli
   ```

2. Example Usage:
   You can read or study the source code in `example/main.go` for an example of its usage.

3. Creating a CLI Application:
   ```go
   package main

   import (
       "fmt"
       gocli "github.com/budimanlai/go-cli"
       "github.com/budimanlai/go-cli/example/commands"
   )

   func main() {
       cli := gocli.NewCliWithConfig(gocli.CliOptions{
           AppName:     "Go CLI Demo",
           Version:     "1.0.0",
           ConfigFile:  []string{"config/main.conf"},
           RuntimePath: "runtime/",
       })

       cli.AddCommand("random_string", commands.RandomString)
       cli.AddCommand("world", commands.World)
       cli.AddCommand("clean_log", commands.Clearlog)
       cli.StartService("run", "mulai", commands.Listen)
       cli.StopService("berhenti")

       e := cli.Run()
       if e != nil {
           fmt.Println(e.Error())
       }
   }
   ```

For more information, please refer to the complete documentation at [documentation link].

## GitHub Link
For more information and the complete source code, visit the [GitHub Repository](https://github.com/budimanlai/go-cli).
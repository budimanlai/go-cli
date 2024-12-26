package commands

import (
	"fmt"
	"time"

	gocli "github.com/budimanlai/go-cli"
)

func Clearlog(c *gocli.Cli) {
	fmt.Println("Simulate Clear log")

	table := []string{"table1_log", "table2_log", "table3_log"}
	for _, element := range table {
		fmt.Println("Clear log table:", element)
		time.Sleep(8 * time.Second)
	}

	fmt.Println("Done...")
}

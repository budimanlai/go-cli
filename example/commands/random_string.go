package commands

import (
	"fmt"
	"math/rand"

	gocli "github.com/budimanlai/go-cli"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func random_string(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func RandomString(c *gocli.Cli) {
	length := c.Args.GetIntOr("length", 10)
	fmt.Println("Generate random string. Default length: 10")
	fmt.Println("Length:", length)
	fmt.Println("String:", random_string(length))
}

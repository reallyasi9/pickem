package main

import (
	"fmt"
	"os"
)

var commands map[string]func([]string) error

func init() {
	commands = make(map[string]func([]string) error)
	commands["teams"] = teams
}

func printUsage() {
	fmt.Println("Usage: download CMD [options...]")
	fmt.Println("")
	fmt.Println("CMD must be one of the following:")
	for key := range commands {
		fmt.Printf("\t%s\n", key)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}
	cmd, args := args[0], args[1:]
	if f, ok := commands[cmd]; ok {
		err := f(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		os.Exit(0)
	}
	fmt.Printf("Command '%s' not recognized\n", cmd)
	printUsage()
	os.Exit(1)
}

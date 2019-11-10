package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
)

var commands map[string]func(context.Context, []string) error = make(map[string]func(context.Context, []string) error)

const apiURL = "https://api.collegefootballdata.com/teams"

var client http.Client

var fs *firestore.Client

func printUsage() {
	fmt.Println("Usage: download CMD [options...]")
	fmt.Println("")
	fmt.Println("CMD must be one of the following:")
	for key := range commands {
		fmt.Printf("\t%s\n", key)
	}
}

func init() {
	client = http.Client{
		// default
	}
	var err error
	fs, err = firestore.NewClient(context.Background(), os.Getenv("GCP_PROJECT"))
	if err != nil {
		panic(err)
	}
}

func main() {
	ctx := context.Background()
	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}
	cmd, args := args[0], args[1:]
	if f, ok := commands[cmd]; ok {
		err := f(ctx, args)
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

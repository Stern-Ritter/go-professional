package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatalf("invalid number of arguments: %d", len(args))
	}

	dir := args[0]
	env, err := ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	cmd := args[1:]
	RunCmd(cmd, env)
}

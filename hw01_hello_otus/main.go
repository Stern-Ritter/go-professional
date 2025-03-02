package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	prompt := reverse.String("Hello, OTUS!")
	fmt.Println(prompt)
}

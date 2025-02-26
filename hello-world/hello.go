package main

import (
	"fmt"

	"rsc.io/quote"

	"golang.org/x/example/hello/reverse"
)

func main() {
	fmt.Println(quote.Hello())
	fmt.Println(reverse.String("Hello"), reverse.Int(24601))
}

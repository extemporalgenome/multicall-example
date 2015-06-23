package main

import (
	"fmt"
	"os"
)

func main() { Main(os.Args[1:]) }

func Main(args []string) {
	fmt.Println("log invoked with args:", args)
	fmt.Println("(outputting log...)")
}

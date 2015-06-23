package main

import (
	"fmt"
	"os"
)

func main() { Main(os.Args[1:]) }

func Main(args []string) {
	fmt.Println("commit invoked with args:", args)
}

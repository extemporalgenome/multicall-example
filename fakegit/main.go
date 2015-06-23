package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cmdcommit "github.com/extemporalgenome/multicall-example/fakegit-commit"
	cmdlog "github.com/extemporalgenome/multicall-example/fakegit-log"
)

const prefix = "fakegit"

var commands = map[string]func([]string){
	"log":    cmdlog.Main,
	"commit": cmdcommit.Main,
}

func exitf(code int, format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(code)
}

func help() {
	for name := range commands {
		fmt.Println(name)
	}
}

func main() {
	args := os.Args
	// check the invoked program name (multicall technique)
	name := filepath.Base(args[0])
	args = args[1:]

	if name != prefix {
		// we've been invoked with a custom os.Args[0] (usually via symlink, hardlink, or syscall)

		const dashprefix = prefix + "-"
		if strings.HasPrefix(name, dashprefix) {
			// trim off the prefix
			i := len(dashprefix)
			name = name[i:]
		}
	} else {
		// this invocation, the program name is literally the prefix, so dispatch based on os.Args[1] instead.

		if len(os.Args) == 1 {
			// no subcommand provided, so list commands
			help()
			return
		}

		name, args = args[0], args[1:]
	}

	// resolve the santized command name to a function we can call
	cmd, ok := commands[name]
	if !ok {
		exitf(2, "unrecognized command %q\n", name)
	}

	cmd(args)
}

package main

import (
	"bufio"
	"fmt"
	"github.com/minero/minero-go/proto/nbt"
	"os"
	"strings"
	. "github.com/logrusorgru/aurora"
)

// Starts the read-evaluate-print-loop
func startRepl(startAt *nbt.Compound) {
	root = startAt
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("[ %s ] >> ", Magenta(curPath))
		if !scanner.Scan() {
			break
		}

		// read
		str := scanner.Text()

		// evaluate
		split := strings.SplitN(str, " ", 2)

		var cmdString, args string
		if len(split) == 1 {
			cmdString = split[0]
			args = ""
		} else {
			cmdString, args = split[0], split[1]
		}

		cmd := commands[cmdString]

		// print
		if cmd == nil {
			fmt.Println(Red("Command not found:"), Red(cmdString))
			continue
		}

		err := cmd(args)
		if err != nil {
			fmt.Println(Red("Error:"), Red(err))
			continue
		}

		// loop (restart loop)
	}
}

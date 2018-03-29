package main

import (
	"bufio"
	"fmt"
	"github.com/minero/minero-go/proto/nbt"
	"os"
	"strings"
)

// Starts the read-evaluate-print-loop
func startRepl(startAt *nbt.Compound) {
	root = startAt
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("%q >>> ", curPath)
	for scanner.Scan() {
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
			fmt.Println("Command not found: " + cmdString)
			continue
		}

		err := cmd(args)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// loop (restart loop)
	}
}

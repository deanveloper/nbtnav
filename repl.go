package main

import (
    "github.com/minero/minero-go/proto/nbt"
    "bufio"
    "os"
    "fmt"
    "strings"
)

func startRepl(startAt *nbt.Compound) {
    root = startAt
    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        fmt.Printf("%q >>> ", curPath)
        str := scanner.Text()

        split := strings.SplitN(str, " ", 2)

        var cmdString, args string
        if len(split) == 1 {
            cmdString = split[0]
            args = ""
        } else {
            cmdString, args = split[0], split[1]
        }

        cmd := commands[cmdString]
        if cmd == nil {
            fmt.Println("Command not found: " + cmdString)
            continue
        }

        err := cmd(args)
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }
    }
}
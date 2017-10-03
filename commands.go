package main

import (
    "github.com/minero/minero-go/proto/nbt"
    "errors"
    "fmt"
    "os"
)

// represents a command
type command func(args string) error

// represents the root of the compound
var root *nbt.Compound

// represents the nbt path we are at
var curPath = "/"

var errNotFound = errors.New("cannot find anything with that path")
var errNotCompound = errors.New("cannot navigate into a non-compound")

// represents a map of command names to the functions they run
var commands map[string]command = map[string]command{
    "cd": cdCommand,
    "ls": lsCommand,
    "exit": exitCommand,
}

// Enters an nbt compound
func cdCommand(args string) error {
    next, err := customLookup(args)
    if err != nil {
        return err
    }
    _, ok := next.(*nbt.Compound)
    if !ok {
        return errNotCompound
    }

    curPath = resolve(curPath, args)

    fmt.Println("Entered", args)

    return nil
}

// View everything inside the current compound
func lsCommand(args string) error {
    if len(args) == 0 {

        tag, _ := customLookup(".")
        prettyPrint(tag.(*nbt.Compound).Value)

    } else {

        path := resolve(curPath, args)
        tag, _ := customLookup(path)

        if comp, ok := tag.(*nbt.Compound); ok {
            prettyPrint(comp.Value)
        } else {
            return errNotCompound
        }
    }

    return nil
}

// Exit the repl
func exitCommand(args string) error {
    os.Exit(0)
    return nil
}
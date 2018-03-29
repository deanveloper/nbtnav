package main

import (
	"errors"
	"fmt"
	"github.com/minero/minero-go/proto/nbt"
	"os"
)

// represents a command
type command func(args string) error

// represents the root of the compound
var root *nbt.Compound

// represents the nbt path we are at
var curPath = "/"

var errNotFound = errors.New("cannot find anything with that path")
var errNotCompound = errors.New("not a compound")
var errIsCompound = errors.New("cannot print out a compound")
var errNotEnoughArgs = errors.New("not enough arguments")

// represents a map of command names to the functions they run
var commands = map[string]command{
	"cd":   cdCommand,
	"ls":   lsCommand,
	"tree": treeCommand,
	"cat":  catCommand,
	"exit": exitCommand,
}

// Enters an nbt compound
func cdCommand(args string) error {
	if args == "--help" {
		fmt.Println("cd <path>")
		return nil
	}

	next, err := nextTag(args)
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
	if args == "--help" {
		fmt.Println("ls [path]")
		return nil
	}

	if args == "" {

		tag, _ := nextTag(".")
		if comp, ok := tag.(*nbt.Compound); ok {
			prettyPrint(comp.Value)
		} else {
			return errNotCompound
		}

	} else {

		path := resolve(curPath, args)
		tag, _ := nextTag(path)

		if comp, ok := tag.(*nbt.Compound); ok {
			prettyPrint(comp.Value)
		} else {
			return errNotCompound
		}
	}

	return nil
}

// Similar to ls, but views the whole tree
func treeCommand(args string) error {
	if args == "--help" {
		fmt.Println("tree [path]")
		return nil
	}

	if args == "" {

		tag, _ := nextTag(".")
		deepPrettyPrint(tag.(*nbt.Compound).Value)

	} else {

		path := resolve(curPath, args)
		tag, _ := nextTag(path)

		if comp, ok := tag.(*nbt.Compound); ok {
			deepPrettyPrint(comp.Value)
		} else {
			return errNotCompound
		}
	}

	return nil
}

// Prints out a value
func catCommand(args string) error {
	if args == "--help" {
		fmt.Println("cat [path]")
		return nil
	}

	if args == "" {
		return errNotEnoughArgs
	} else {

		path := resolve(curPath, args)
		tag, _ := nextTag(path)

		if _, ok := tag.(*nbt.Compound); ok {
			return errIsCompound
		} else {
			fmt.Println(tag)
		}
	}

	return nil
}

// Exit the repl
func exitCommand(args string) error {
	os.Exit(0)
	return nil
}

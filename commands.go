package main

import (
	"fmt"
	"github.com/minero/minero-go/proto/nbt"
	"os"
	"strings"
	"path"
	"sort"
	"strconv"
	. "github.com/logrusorgru/aurora"
	"bufio"
	"io"
	"compress/gzip"
	"compress/zlib"
)

// represents a command
type command func(args string) error

// represents the root of the compound
var root *nbt.Compound

// represents the nbt path we are at
var curPath = "/"

// represents a map of command names to the functions they run
var commands = map[string]command{
	"help": helpCommand,
	"cd":   cdCommand,
	"ls":   lsCommand,
	"tree": treeCommand,
	"cat":  catCommand,
	"set":  setCommand,
	"save": saveCommand,
	"exit": exitCommand,
}

var help = map[string]string{
	"help":                     "Shows command list.",
	"cd <compound>":            "Switches context to the provided compound. \"..\" supported.",
	"ls [compound]":            "Lists the tags in the current context, or provided compound.",
	"tree [compound]":          "Same as ls, but recursive.",
	"cat <tag>":                "Prints out the value of the provided tag.",
	"set <tag> <type> [value]": "Sets a tag's value/type.",
	"save [compress] [output]": "Saves to output. Compression can be gzip, zlib, or none.",
	"exit":                     "Exits nbtnav.",
}

func helpCommand(arg string) error {
	keys := make([]string, len(help))
	i := 0
	maxLen := -1
	for k := range help {
		keys[i] = k
		i++
		if len(k) > maxLen {
			maxLen = len(k)
		}
	}
	max := strconv.Itoa(maxLen)
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%"+max+"s - %s\n", Blue(k), Cyan(help[k]))
	}

	return nil
}

// Enters an nbt compound
func cdCommand(arg string) error {

	next, err := pathToTag(arg)
	if err != nil {
		return err
	}
	_, ok := next.(*nbt.Compound)
	if !ok {
		return errNotCompound
	}

	curPath = resolve(curPath, arg)

	return nil
}

// View everything inside the current compound
func lsCommand(arg string) error {

	if arg == "" {

		tag, _ := pathToTag(".")
		if comp, ok := tag.(*nbt.Compound); ok {
			prettyPrint(comp.Value)
		} else {
			return errNotCompound
		}

	} else {

		next := resolve(curPath, arg)
		tag, _ := pathToTag(next)

		if comp, ok := tag.(*nbt.Compound); ok {
			prettyPrint(comp.Value)
		} else {
			return errNotCompound
		}
	}

	return nil
}

// Similar to ls, but views the whole tree
func treeCommand(arg string) error {

	if arg == "" {

		tag, _ := pathToTag(".")
		deepPrettyPrint(tag.(*nbt.Compound).Value)

	} else {

		next := resolve(curPath, arg)
		tag, _ := pathToTag(next)

		if comp, ok := tag.(*nbt.Compound); ok {
			deepPrettyPrint(comp.Value)
		} else {
			return errNotCompound
		}
	}

	return nil
}

// Prints out a value
func catCommand(arg string) error {

	if arg == "" {
		return errNotEnoughArgs
	} else {

		next := resolve(curPath, arg)
		tag, err := pathToTag(next)
		if err != nil {
			return err
		}

		_, ok32 := tag.(*nbt.Float32)
		_, ok64 := tag.(*nbt.Float64)
		if ok32 || ok64 {
			fmt.Println(prettyFloat(tag, true))
		}
		if _, ok := tag.(*nbt.Compound); ok {
			return errPrintedCompound
		} else if barr, ok := tag.(*nbt.ByteArray); ok {
			fmt.Println(prettyByteArray(barr, true))
		} else {
			fmt.Println(prettyString(tag))
		}
	}

	return nil
}

// Exit the repl
func exitCommand(arg string) error {
	os.Exit(0)
	return nil
}

func setCommand(arg string) error {
	args := parseMultiArgs(arg)

	if len(args) < 2 {
		return errNotEnoughArgs
	}

	targ := args[0]

	typ, err := typeFromString(args[1])
	if err != nil {
		return err
	}

	val := ""
	if len(args) >= 3 {
		val = args[2]
	}

	parentPath, tagName := path.Split(resolve(curPath, targ))
	if tagName == "" {
		return errNotFound
	}

	newTag := typ.New()
	setTagValue(newTag, val)

	parentTag, err := pathToTag(parentPath)
	if err != nil {
		return err
	}

	if comp, ok := parentTag.(*nbt.Compound); ok {
		comp.Value[tagName] = newTag
		return nil
	} else {
		return errNotFound
	}
}

func saveCommand(arg string) error {
	args := parseMultiArgs(arg)
	var compress, output string

	if len(args) >= 1 {
		compress = strings.ToLower(args[0])
	} else {
		compress = "none"
	}
	if len(args) >= 2 {
		output = args[1]
	} else {
		output = os.Args[1]
	}

	if compress != "gzip" && compress != "zlib" && compress != "none" {
		return errInvalidCompression
	}

	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	var w io.Writer = f

	if compress == "gzip" {
		gWriter := gzip.NewWriter(w)
		defer gWriter.Close()
		w = gWriter
	} else if compress == "zlib" {
		zWriter := zlib.NewWriter(w)
		defer zWriter.Close()
		w = zWriter
	}

	err = nbt.Write(w, root)
	if err != nil {
		return err
	}

	return nil
}

// Simple way to parse multiple arguments into a slice. Not very advanced but avoids
// needing to use some giant framework.
func parseMultiArgs(args string) []string {
	scan := bufio.NewScanner(strings.NewReader(args))
	scan.Split(bufio.ScanWords)

	stringMode := false
	var argSlice []string
	index := -1

	for scan.Scan() {
		txt := scan.Text()

		if !stringMode {
			index++
			argSlice = append(argSlice, "")
		} else {
			argSlice[index] += " "
		}

		oldTxt := txt
		if txt = strings.TrimPrefix(oldTxt, `"`); oldTxt != txt {
			stringMode = true
		}
		oldTxt = txt
		if txt = strings.TrimSuffix(oldTxt, `"`); oldTxt != txt {
			stringMode = false
		}

		argSlice[index] += txt
	}

	return argSlice
}

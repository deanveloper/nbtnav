package main

import (
	"fmt"
	"github.com/minero/minero-go/proto/nbt"
	"sort"
	"path"
	"strings"
	. "github.com/logrusorgru/aurora"
)

// Essentially path.Join but will also clean.
func resolve(from, to string) string {
	if strings.HasPrefix(to, "/") {
		return path.Clean(to)
	}
	return path.Clean(path.Join(from, to))
}

// Gets a tag with a given path, relative to the current tag
func nextTag(nbtPath string) (nbt.Tag, error) {
	absPath := resolve(curPath, nbtPath)

	if absPath == "/" {
		return root, nil
	}

	// Remove leading slash
	absPath = absPath[1:]

	split := strings.SplitN(absPath, "/", 2)
	next := root.Value[split[0]]
	for len(split) == 1 || split[1] != "" {
		split = strings.SplitN(absPath, "/", 2)
		next = root.Value[split[0]]

		if next == nil {
			return nil, errNotFound
		}
	}
	return next, nil
}

// Pretty-Print
func prettyPrint(tags map[string]nbt.Tag) {
	// sort keys
	keys := make([]string, len(tags))
	i := 0
	for key := range tags {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := tags[key]
		if comp, ok := value.(*nbt.Compound); ok {
			fmt.Printf("%s: (%s len(%d))\n", Blue(key), Green(comp.Type()), Blue(len(comp.Value)))
		} else if list, ok := value.(*nbt.List); ok {
			fmt.Printf("%s: (%s len(%d)) %s\n", Blue(key), Green(list.Type()), Blue(len(list.Value)), Cyan(list.Value))
		} else {
			fmt.Printf("%s: %s\n", Blue(key), Cyan(value))
		}
	}
}

// Ease-of-use entry method for deepPrettyPrintRecur
func deepPrettyPrint(nbt map[string]nbt.Tag) {
	deepPrettyPrintRecur(0, nbt)
}

// Recursively prints an nbt tree with a beautiful tree structure.
func deepPrettyPrintRecur(deepness int, tags map[string]nbt.Tag) {
	// sort keys
	keys := make([]string, len(tags))
	i := 0
	for key := range tags {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	// https://en.wikipedia.org/wiki/Box-drawing_character
	// Using code points U+2500 (─), U+2502 (│), U+251C (├), and U+2514 (└)
	for index, key := range keys {

		prefix := ""
		for i := 0; i < deepness; i++ {
			prefix += "│   "
		}
		if index == len(keys)-1 {
			prefix += "└───"
		} else {
			prefix += "├───"
		}

		value := tags[key]
		if comp, ok := value.(*nbt.Compound); ok {
			fmt.Printf("%s %s: (%s len(%d))\n", prefix, Blue(key), Green(comp.Type()), Blue(len(comp.Value)))
			deepPrettyPrintRecur(deepness+1, comp.Value)
		} else if list, ok := value.(*nbt.List); ok {
			fmt.Printf("%s %s: (%s len(%d)) %s\n", prefix, Blue(key), Green(list.Type()), Blue(len(list.Value)), Cyan(list.Value))
		} else {
			fmt.Printf("%s %s: %s\n", prefix, Blue(key), Cyan(tags[key]))
		}
	}
}

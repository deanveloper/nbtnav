package main

import (
	"fmt"
	"github.com/minero/minero-go/proto/nbt"
	"sort"
	"path"
	"strings"
)

// Essentially path.Join but will also clean.
func resolve(from, to string) string {
	if strings.HasPrefix(to, "/") {
		return path.Clean(to)
	}
	return path.Clean(path.Join(from, to))
}

// Gets a tag with a given path
func nextTag(path string) (nbt.Tag, error) {
	nextPath := resolve(curPath, path)

	next := root.Value[nextPath[1:]]
	if next == nil {
		return nil, errNotFound
	}
	return next, nil
}

// Pretty-Print
func prettyPrint(tags map[string]nbt.Tag) {
	// sort keys
	keys := make([]string, len(tags))
	sort.Strings(keys)

	for _, key := range keys {
		value := tags[key]
		if comp, ok := value.(*nbt.Compound); ok {
			fmt.Printf("%q: (%s len(%d))", key, comp.Type(), len(comp.Value))
		} else {
			fmt.Printf("%q: %s", key, tags[key])
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
			fmt.Printf("%s%q: (%s len(%d))", prefix, key, comp.Type(), len(comp.Value))
			deepPrettyPrintRecur(deepness+1, comp.Value)
		} else {
			fmt.Printf("%s%q: %s", prefix, key, tags[key])
		}
	}
}

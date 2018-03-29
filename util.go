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
		fmt.Printf("%s: %s\n", Blue(key), prettyString(tags[key]))
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

		fmt.Printf("%s %s: %s\n", prefix, Blue(key), prettyString(tags[key]))
	}
}

func prettyString(tag nbt.Tag) string {
	if comp, ok := tag.(*nbt.Compound); ok {
		return fmt.Sprintf("(%s len(%d))", Green(comp.Type()), Blue(len(comp.Value)))
	} else if list, ok := tag.(*nbt.List); ok {
		return fmt.Sprintf("(%s len(%d))", Green(list.Type()), Blue(len(list.Value)))
	} else {
		return fmt.Sprintf("(%s) %s", Green(tag.Type()), Blue(tag))
	}
}

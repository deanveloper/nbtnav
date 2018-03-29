package main

import (
	"fmt"
	"github.com/minero/minero-go/proto/nbt"
	"sort"
	"path"
	"strings"
	. "github.com/logrusorgru/aurora"
	"encoding/hex"
	"bytes"
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
	for len(split) > 1 && split[1] != "" {
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

		if comp, ok := tags[key].(*nbt.Compound); ok {
			deepPrettyPrintRecur(deepness + 1, comp.Value)
		}
	}
}

// Prints a tag out in a better-looking way
func prettyString(tag nbt.Tag) string {
	// Byte
	if v, ok := tag.(*nbt.Int8); ok {
		return fmt.Sprintf("(%s) %s", Green(tag.Type().String()[3:]), Cyan(v.Int8))
	}
	// Short
	if v, ok := tag.(*nbt.Int16); ok {
		return fmt.Sprintf("(%s) %s", Green(tag.Type().String()[3:]), Cyan(v.Int16))
	}
	// Int
	if v, ok := tag.(*nbt.Int32); ok {
		return fmt.Sprintf("(%s) %s", Green(tag.Type().String()[3:]), Cyan(v.Int32))
	}
	// Long
	if v, ok := tag.(*nbt.Int64); ok {
		return fmt.Sprintf("(%s) %s", Green(tag.Type().String()[3:]), Cyan(v.Int64))
	}
	// Float
	if v, ok := tag.(*nbt.Float32); ok {
		return fmt.Sprintf("(%s) %s", Green(tag.Type().String()[3:]), Cyan(v.Float32))
	}
	// Double
	if v, ok := tag.(*nbt.Float64); ok {
		return fmt.Sprintf("(%s) %s", Green(tag.Type().String()[3:]), Cyan(v.Float64))
	}
	// ByteArray
	if v, ok := tag.(*nbt.ByteArray); ok {
		return prettyByteArray(v, false)
	}
	// String
	if v, ok := tag.(*nbt.String); ok {
		return fmt.Sprintf("(%s) %s", Green(tag.Type().String()[3:]), Cyan(v.Value))
	}
	// List
	if v, ok := tag.(*nbt.List); ok {
		return fmt.Sprintf("(%s(%s) len(%d))", Green(tag.Type().String()[3:]), Green(v.Typ.String()[3:]), Blue(len(v.Value)))
	}
	// Compound
	if v, ok := tag.(*nbt.Compound); ok {
		return fmt.Sprintf("(%s len(%d))", Green(tag.Type().String()[3:]), Blue(len(v.Value)))
	}
	return fmt.Sprintf("(%s) %s", Green(tag.Type().String()[3:]), Cyan(tag))
}

// Function specifically for printing out byte arrays
func prettyByteArray(tag *nbt.ByteArray, longForm bool) string {
	var buf bytes.Buffer

	if longForm || len(tag.Value) <= 40 {
		for i := 0; i < len(tag.Value); i++ {
			buf.WriteByte(byte(tag.Value[i]))
		}
	} else {
		for i := 0; i < 37; i++ {
			buf.WriteByte(byte(tag.Value[i]))
		}
	}

	str := hex.EncodeToString(buf.Bytes())
	if len(str) >= 40 {
		str = str[:37]
		str += "..."
	}
	return fmt.Sprintf("(%s len(%d)) %s", Green(tag.Type().String()[3:]), Blue(len(tag.Value)), Cyan(str))
}
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
	"github.com/minero/minero/types"
)

// Essentially path.Join but will also clean.
func resolve(from, to string) string {
	if strings.HasPrefix(to, "/") {
		return path.Clean(to)
	}
	return path.Clean(path.Join(from, to))
}

// Gets a tag with a given path, relative to the current tag
func pathToTag(nbtPath string) (nbt.Tag, error) {
	absPath := resolve(curPath, nbtPath)

	if absPath == "/" {
		return root, nil
	}

	// Remove leading slash
	absPath = absPath[1:]

	next := root.Lookup(absPath)
	if next == nil {
		return nil, errNotFound
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
func deepPrettyPrintRecur(depth int, tags map[string]nbt.Tag) {
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
		for i := 0; i < depth; i++ {
			prefix += "│   "
		}
		if index == len(keys)-1 {
			prefix += "└───"
		} else {
			prefix += "├───"
		}

		fmt.Printf("%s %s: %s\n", prefix, Blue(key), prettyString(tags[key]))

		if comp, ok := tags[key].(*nbt.Compound); ok {
			deepPrettyPrintRecur(depth + 1, comp.Value)
		}
	}
}

// Prints a tag out in a better-looking way
func prettyString(tag nbt.Tag) string {
	// Byte
	if v, ok := tag.(*nbt.Int8); ok {
		return fmt.Sprintf("(%s) %d", Green(tag.Type().String()[3:]), Cyan(v.Int8))
	}
	// Short
	if v, ok := tag.(*nbt.Int16); ok {
		return fmt.Sprintf("(%s) %d", Green(tag.Type().String()[3:]), Cyan(v.Int16))
	}
	// Int
	if v, ok := tag.(*nbt.Int32); ok {
		return fmt.Sprintf("(%s) %d", Green(tag.Type().String()[3:]), Cyan(v.Int32))
	}
	// Long
	if v, ok := tag.(*nbt.Int64); ok {
		return fmt.Sprintf("(%s) %d", Green(tag.Type().String()[3:]), Cyan(v.Int64))
	}
	// Float/Double
	_, ok32 := tag.(*nbt.Float32)
	_, ok64 := tag.(*nbt.Float64)
	if ok32 || ok64 {
		return prettyFloat(tag, false)
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

// Function specifically for printing out floats. Panics if tag is not a float.
func prettyFloat(tag nbt.Tag, longForm bool) string {
	if longForm {
		// Float
		if v, ok := tag.(*nbt.Float32); ok {
			return fmt.Sprintf("(%s) %f", Green(tag.Type().String()[3:]), Cyan(v.Float32))
		}
		// Double
		if v, ok := tag.(*nbt.Float64); ok {
			return fmt.Sprintf("(%s) %f", Green(tag.Type().String()[3:]), Cyan(v.Float64))
		}
	} else {
		// Float
		if v, ok := tag.(*nbt.Float32); ok {
			return fmt.Sprintf("(%s) %.2f", Green(tag.Type().String()[3:]), Cyan(v.Float32))
		}
		// Double
		if v, ok := tag.(*nbt.Float64); ok {
			return fmt.Sprintf("(%s) %.2f", Green(tag.Type().String()[3:]), Cyan(v.Float64))
		}
	}

	panic("Provided tag was not a float type")
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

// Function to set a tag's value.
// Automatically parses string, returns error if not able to parse.
//
// Cannot set value of compounds and lists
func setTagValue(tag nbt.Tag, val string) error {
	// Byte
	if v, ok := tag.(*nbt.Int8); ok {
		_, err := fmt.Sscan(val, &v.Int8)
		return err
	}
	// Short
	if v, ok := tag.(*nbt.Int16); ok {
		_, err := fmt.Sscan(val, &v.Int16)
		return err
	}
	// Int
	if v, ok := tag.(*nbt.Int32); ok {
		_, err := fmt.Sscan(val, &v.Int32)
		return err
	}
	// Long
	if v, ok := tag.(*nbt.Int64); ok {
		_, err := fmt.Sscan(val, &v.Int64)
		return err
	}
	// Float
	if v, ok := tag.(*nbt.Float32); ok {
		_, err := fmt.Sscan(val, &v.Float32)
		return err
	}
	// Double
	if v, ok := tag.(*nbt.Float64); ok {
		_, err := fmt.Sscan(val, &v.Float64)
		return err
	}
	// ByteArray
	if v, ok := tag.(*nbt.ByteArray); ok {
		slice, err := hex.DecodeString(val)
		if err != nil {
			return err
		}
		int8Slice := make([]types.Int8, len(slice))
		for i := 0; i < len(slice); i++ {
			int8Slice[i] = types.Int8(slice[i])
		}
		v.Value = int8Slice
	}
	// String
	if v, ok := tag.(*nbt.String); ok {
		v.Value = val
	}

	return nil
}

func typeFromString(str string) (nbt.TagType, error) {
	for tag := nbt.TagEnd; tag <= nbt.TagIntArray; tag++ {
		if strings.ToLower(str) == strings.ToLower(tag.String()[3:]) {
			return tag, nil
		}
	}

	return -1, errInvalidTagType
}

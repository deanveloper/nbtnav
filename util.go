package main

import (
    "strings"
    "github.com/minero/minero-go/proto/nbt"
    "fmt"
    "sort"
)

// Resolves an NBT path
func resolve(from, to string) string {
    split := strings.SplitN(to, "/", 2)
    nextPath := split[0]
    if nextPath == ".." {
        from = from[:len(curPath) - 1][:strings.LastIndex(from, "/") + 1]
    } else if nextPath == "." {
        // do nothing
    } else if nextPath == "" {
        from = "/"
    } else {
        from = curPath + nextPath
    }

    if from == "" {
        from = "/"
    }

    if len(split) == 2 {
        return resolve(from, split[1])
    }
    return from
}

// Performs a lookup on a compound. This is not a command.
func customLookup(arg string) (nbt.Tag, error) {
    nextPath := resolve(curPath, arg)

    if nextPath == "" {
        return nil, errNotFound
    }

    next := root.Value[nextPath[1:]]
    if next == nil {
        return nil, errNotFound
    }

    return next, nil
}

// Pretty-Print
func prettyPrint(nbt map[string]nbt.Tag) {
    // sort keys
    keys := make([]string, len(nbt))
    sort.Strings(keys)

    for _, key := range keys {
        if comp, ok := nbt[key].(*nbt.Compound); ok {
            fmt.Printf("%q [%s len(%d)]", key, comp.Type(), len(comp.Value))
        } else {
            fmt.Printf("%q [%s]", key, nbt[key].Type())
        }
    }
}

// Ease-of-use entry method for deepPrettyPrintRecur
func deepPrettyPrint(nbt map[string]nbt.Tag) {
    deepPrettyPrintRecur(0, nbt)
}

// Recursively prints an nbt tree with a beautiful tree structure.
func deepPrettyPrintRecur(deepness int, nbt map[string]nbt.Tag) {
    // sort keys
    keys := make([]string, len(nbt))
    sort.Strings(keys)

    // https://en.wikipedia.org/wiki/Box-drawing_character
    // Using code points U+2500 (─), U+2502 (│), U+251C (├), and U+2514 (└)
    for index, key := range keys {

        prefix := ""
        for i := 0; i < deepness; i++ {
            prefix += "│   "
        }
        if index == len(keys) - 1 {
            prefix += "└───"
        } else {
            prefix += "├───"
        }

        if comp, ok := nbt[key].(*nbt.Compound); ok {
            fmt.Printf("%s%q [%s len(%d)]", prefix, key, comp.Type(), len(comp.Value))
            deepPrettyPrintRecur(deepness + 1, comp.Value)
        } else {
            fmt.Printf("%s%q [%s]", prefix, key, nbt[key].Type())
        }
    }
}
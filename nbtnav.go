package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/minero/minero-go/proto/nbt"
	"io/ioutil"
	"os"
)

func main() {
	root := getCompoundFromArgs()

	startRepl(root)
}

func getCompoundFromArgs() *nbt.Compound {
	var root *nbt.Compound

	if len(os.Args) == 2 {
		if os.Args[1] == "--help" {
			fmt.Println("nbtnav `filename`")
			os.Exit(0)
		}

		file, err := ioutil.ReadFile(os.Args[1])
		checkErr(err)

		root, err = nbt.Read(bytes.NewReader(file))
		checkErr(err)

	} else {
		checkErr(errors.New("wrong number of arguments, use --help for help"))
	}

	return root
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
}

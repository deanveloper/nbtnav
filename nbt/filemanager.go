package nbt

import (
	"bytes"
	"compress/zlib"
	"io"
	"io/ioutil"
	"compress/gzip"
)

//
func ReadFile(path string) (*Tag, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var nbtReader io.Reader

	tryZlib := bytes.NewReader(data)
	zReader, err := zlib.NewReader(tryZlib)
	if zReader != nil {
		defer zReader.Close()
	}

	if err == nil {
		nbtReader = zReader
	} else {
		// make a new reader, stupid I know
		tryGzip := bytes.NewReader(data)
		gReader, err := gzip.NewReader(tryGzip)
		if gReader != nil {
			defer gReader.Close()
		}

		if err == nil {
			nbtReader = gReader
		} else {
			nbtReader = bytes.NewReader(data)
		}
	}

	// nbtReader is now set
	return Parse(nbtReader)
}

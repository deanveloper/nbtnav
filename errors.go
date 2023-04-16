package main

import "errors"

var errNotFound = errors.New("cannot find tag with that path")
var errNotCompound = errors.New("not a compound")
var errPrintedCompound = errors.New("cannot print out a compound")
var errInvalidTagType = errors.New("not a valid tag type")
var errNotEnoughArgs = errors.New("not enough arguments, use help <cmd> for help")
var errInvalidCompression = errors.New("invalid compression: not gzip, zlib, or none")
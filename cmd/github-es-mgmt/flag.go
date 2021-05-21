package main

import (
	"flag"
	"fmt"
)

type Usager interface {
	Usage()
}

type FlagSet struct {
	*flag.FlagSet
	errorMessage string
}

func NewFlagSet(usage string) *FlagSet {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprint(fs.Output(), usage)
		fs.PrintDefaults()
	}
	return &FlagSet{FlagSet: fs}
}

func (fs *FlagSet) SetError(message string) *FlagSet {
	fs.errorMessage = message
	return fs
}

func (fs *FlagSet) Usage() {
	fs.FlagSet.Usage()
	if fs.errorMessage != "" {
		fmt.Fprintf(fs.Output(), "\n%s\n", fs.errorMessage)
	}
}

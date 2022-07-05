package main

import (
	"flag"
	"fmt"
)

func NewFlagSet(usage string) *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprint(fs.Output(), usage)
		fs.PrintDefaults()
	}
	return fs
}

type UsageError struct {
	fs      *flag.FlagSet
	message string
}

func NewUsageError(fs *flag.FlagSet, message string) *UsageError {
	return &UsageError{fs: fs, message: message}
}

func (e UsageError) Error() string {
	return fmt.Sprintf("usage error: %s", e.message)
}

func (e *UsageError) Usage() {
	e.fs.Usage()
	if e.message != "" {
		fmt.Fprintf(e.fs.Output(), "\n%s\n", e.message)
	}
}

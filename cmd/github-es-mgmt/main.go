package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	version string
	commit  string
	date    string
)

func main() {
	os.Exit(run())
}

func run() int {
	c, fs, err := TopLevelArgsParser{}.Parse(filepath.Base(os.Args[0]), os.Args[1:])
	if fs != nil {
		fs.Usage()
		if err != nil {
			fmt.Fprintf(fs.Output(), "\n%s\n", err)
		}
		return 2
	}
	if err := c.Execute(); err != nil {
		log.Printf("%s", err)
		return 1
	}
	return 0
}

type Command interface {
	Execute() error
}

type TopLevelArgsParser struct{}

func (p TopLevelArgsParser) Parse(command string, args []string) (Command, *flag.FlagSet, error) {
	usage := fmt.Sprintf(`Usage: %s <subcommand> [options]

subcommands:
    certificate    certificate subcommand.
    maintenance    maintenance subcommand.
    version        show version

Run %s <subcommand> -h to show help for subcommand.
`, command, command)
	fs := newFlagSet(nil, usage)
	if err := fs.Parse(args); err != nil {
		return nil, fs, nil
	}
	args = fs.Args()
	if len(args) == 0 {
		return nil, fs, nil
	}

	switch args[0] {
	case "certificate":
		return CertificateArgsParser{}.Parse(command, args[:1:1], args[1:])
	case "maintenance":
		return MaintenanceArgsParser{}.Parse(command, args[:1:1], args[1:])
	case "version":
		return &VersionCommand{}, nil, nil
	default:
		return nil, fs, nil
	}
}

func newFlagSet(subcommands []string, usage string) *flag.FlagSet {
	var name string
	if len(subcommands) > 0 {
		name = strings.Join(subcommands, " ")
	}
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprint(fs.Output(), usage)
		fs.PrintDefaults()
	}
	return fs
}

type VersionCommand struct{}

func (c *VersionCommand) Execute() error {
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Commit:  %s\n", commit)
	fmt.Printf("Date:    %s\n", date)
	return nil
}

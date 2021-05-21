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
	date    string
)

func main() {
	os.Exit(run())
}

func run() int {
	c, usage := TopLevelArgsParser{}.Parse(filepath.Base(os.Args[0]), os.Args[1:])
	if usage != nil {
		usage()
		return 2
	}
	if err := c.Execute(); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}

type Usage func()

type Command interface {
	Execute() error
}

type TopLevelArgsParser struct{}

func (p TopLevelArgsParser) Parse(command string, args []string) (Command, Usage) {
	usage := fmt.Sprintf(`Usage: %s <subcommand> [options]

subcommands:
    certificate    certificate subcommand.
    maintenance    maintenance subcommand.
    version        show version

Run %s <subcommand> -h to show help for subcommand.
`, command, command)
	fs := newFlagSet(nil, usage)
	if err := fs.Parse(args); err != nil {
		return nil, newUsage(fs, "")
	}
	args = fs.Args()
	if len(args) == 0 {
		return nil, newUsage(fs, "")
	}

	switch args[0] {
	case "certificate":
		return CertificateArgsParser{}.Parse(command, args[:1:1], args[1:])
	case "maintenance":
		return MaintenanceArgsParser{}.Parse(command, args[:1:1], args[1:])
	case "version":
		return &VersionCommand{}, nil
	default:
		return nil, newUsage(fs, "")
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

func newUsage(fs *flag.FlagSet, additionalMessage string) Usage {
	return func() {
		fs.Usage()
		if additionalMessage != "" {
			fmt.Fprintf(fs.Output(), "\n%s\n", additionalMessage)
		}
	}
}

type VersionCommand struct{}

func (c *VersionCommand) Execute() error {
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Date:    %s\n", date)
	return nil
}

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	version string
	date    string
)

func main() {
	os.Exit(run())
}

func run() int {
	c, u := TopLevelArgsParser{}.Parse(filepath.Base(os.Args[0]), os.Args[1:])
	if u != nil {
		u.Usage()
		return 2
	}
	if err := c.Execute(); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}

type Command interface {
	Execute() error
}

type TopLevelArgsParser struct{}

func (p TopLevelArgsParser) Parse(command string, args []string) (Command, Usager) {
	usage := fmt.Sprintf(`Usage: %s <subcommand> [options]

subcommands:
    certificate    certificate subcommand.
    maintenance    maintenance subcommand.
    settings       settings subcommand.
    version        show version

Run %s <subcommand> -h to show help for subcommand.

Source repository is https://github.com/hnakamur/github-es-mgmt
`, command, command)
	fs := NewFlagSet(usage)
	if err := fs.Parse(args); err != nil {
		return nil, fs
	}
	args = fs.Args()
	if len(args) == 0 {
		return nil, fs
	}

	switch args[0] {
	case "certificate":
		return CertificateArgsParser{}.Parse(command, args[:1:1], args[1:])
	case "maintenance":
		return MaintenanceArgsParser{}.Parse(command, args[:1:1], args[1:])
	case "settings":
		return SettingsArgsParser{}.Parse(command, args[:1:1], args[1:])
	case "version":
		return &VersionCommand{}, nil
	default:
		return nil, fs
	}
}

type VersionCommand struct{}

func (c *VersionCommand) Execute() error {
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Date:    %s\n", date)
	return nil
}

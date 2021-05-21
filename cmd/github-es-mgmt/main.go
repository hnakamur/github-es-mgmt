package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const globalUsage = `Usage: %s <subcommand> [options]

subcommands:
  cert set              Set certificate.
  maintenance status    Get maintenance status.
  maintenance enable    Enable maintenance mode.
  maintenance disable   Disable maintenance mode.
  version               Show version

Run %s <subcommand> -h to show help for subcommand.
`

var cmdName = filepath.Base(os.Args[0])

var (
	version string
	commit  string
	date    string
)

func main() {
	os.Exit(run())
}

func run() int {
	flag.Usage = func() {
		fmt.Printf(globalUsage, cmdName, cmdName)
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return 2
	}

	var c Command
	var subcommandLen int
	switch args[0] {
	case "cert":
		if len(args) == 1 || args[1] != "set" {
			flag.Usage()
			return 2
		}
		c = &CertSetCommand{}
	case "maintenance":
		if len(args) == 1 {
			flag.Usage()
			return 2
		}
		switch args[1] {
		case "status":
			c = &MaintenanceStatusCommand{}
			subcommandLen = 2
		case "enable":
			c = &MaintenanceSetCommand{enabled: true}
			subcommandLen = 2
		case "disable":
			c = &MaintenanceSetCommand{enabled: false}
			subcommandLen = 2
		}
	default:
		flag.Usage()
		return 2
	}

	subcommands := args[:subcommandLen]
	fs := buildSubdommandFlagSet(c, subcommands)
	if err := c.Parse(fs, args[subcommandLen:]); err != nil {
		log.Printf("%s", err)
		return 2
	}
	if err := c.Execute(); err != nil {
		log.Printf("%s", err)
		return 1
	}
	return 0
}

func runShowVersion(args []string) error {
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Commit:  %s\n", commit)
	fmt.Printf("Date:    %s\n", date)
	return nil
}

type Command interface {
	UsageTemplate() string
	Parse(fs *flag.FlagSet, args []string) error
	Execute() error
}

func buildSubdommandFlagSet(c Command, subcommands []string) *flag.FlagSet {
	name := strings.Join(subcommands, " ")
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	fs.Usage = func() {
		usageStr := strings.ReplaceAll(c.UsageTemplate(), "{{command}}", cmdName)
		fmt.Fprintf(fs.Output(), "%s", usageStr)
		fs.PrintDefaults()
	}
	return fs
}

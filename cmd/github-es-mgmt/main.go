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
  set-cert            Set certificate.
  get-maintenance     Get maintenance status.
  set-maintenance     Enable or disable maintenance mode.
  version             Show version

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
	switch args[0] {
	case "set-cert":
		c = &SetCertCommand{}
	case "get-maintenance":
		c = &GetMaintenanceCommand{}
	case "set-maintenance":
		c = &SetMaintenanceCommand{}
	default:
		flag.Usage()
		return 2
	}

	fs := buildSubdommandFlagSet(c, args)
	if err := c.Parse(fs, args[1:]); err != nil {
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

func buildSubdommandFlagSet(c Command, args []string) *flag.FlagSet {
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	fs.Usage = func() {
		usageStr := strings.ReplaceAll(c.UsageTemplate(), "{{command}}", cmdName)
		fmt.Fprintf(fs.Output(), "%s", usageStr)
		fs.PrintDefaults()
	}
	return fs
}

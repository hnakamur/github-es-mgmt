package main

import (
	"flag"
	"fmt"
	"strings"
)

type MaintenanceArgsParser struct{}

func (p MaintenanceArgsParser) Parse(command string, subcommands, args []string) (Command, *flag.FlagSet, error) {
	usageTemplate := fmt.Sprintf(`Usage: %s %s <subcommand> [options]

subcommands:
    status    Get maintenance status.
    enable    Enable maintenance mode.
    disable   Disable maintenance mode.
`, command, strings.Join(subcommands, " "))
	fs := newFlagSet(subcommands, usageTemplate)
	if err := fs.Parse(args); err != nil {
		return nil, fs, nil
	}
	args = fs.Args()
	if len(args) == 0 {
		return nil, fs, nil
	}

	switch args[0] {
	case "status":
		return MaintenanceStatusArgsParser{}.Parse(command, append(subcommands, args[0]), args[1:])
	case "enable":
		return MaintenanceEnableArgsParser{}.Parse(command, append(subcommands, args[0]), args[1:])
	case "disable":
		return MaintenanceDisableArgsParser{}.Parse(command, append(subcommands, args[0]), args[1:])
	default:
		return nil, fs, nil
	}
}

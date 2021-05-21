package main

import (
	"fmt"
	"strings"
)

type MaintenanceArgsParser struct{}

func (p MaintenanceArgsParser) Parse(command string, subcommands, args []string) (Command, Usage) {
	usage := fmt.Sprintf(`Usage: %s %s <subcommand> [options]

subcommands:
    status    Get maintenance status.
    enable    Enable maintenance mode.
    disable   Disable maintenance mode.
`, command, strings.Join(subcommands, " "))
	fs := newFlagSet(subcommands, usage)
	if err := fs.Parse(args); err != nil {
		return nil, newUsage(fs, "")
	}
	args = fs.Args()
	if len(args) == 0 {
		return nil, newUsage(fs, "")
	}

	switch args[0] {
	case "status":
		return MaintenanceStatusArgsParser{}.Parse(command, append(subcommands, args[0]), args[1:])
	case "enable":
		return MaintenanceEnableArgsParser{}.Parse(command, append(subcommands, args[0]), args[1:])
	case "disable":
		return MaintenanceDisableArgsParser{}.Parse(command, append(subcommands, args[0]), args[1:])
	default:
		return nil, newUsage(fs, "")
	}
}

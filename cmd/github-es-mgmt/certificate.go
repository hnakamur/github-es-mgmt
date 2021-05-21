package main

import (
	"flag"
	"fmt"
	"strings"
)

type CertificateArgsParser struct{}

func (p CertificateArgsParser) Parse(command string, subcommands, args []string) (Command, *flag.FlagSet, error) {
	usage := fmt.Sprintf(`Usage: %s %s <subcommand> [options]

subcommands:
    set       Set certificate.
`, command, strings.Join(subcommands, " "))
	fs := newFlagSet(subcommands, usage)
	if err := fs.Parse(args); err != nil {
		return nil, fs, nil
	}
	args = fs.Args()
	if len(args) == 0 {
		return nil, fs, nil
	}

	switch args[0] {
	case "set":
		return CertificateSetArgsParser{}.Parse(command, append(subcommands, args[0]), args[1:])
	default:
		return nil, fs, nil
	}
}

package main

import (
	"fmt"
	"strings"
)

type CertificateArgsParser struct{}

func (p CertificateArgsParser) Parse(command string, subcommands, args []string) (Command, Usager) {
	usage := fmt.Sprintf(`Usage: %s %s <subcommand> [options]

subcommands:
    set       Set certificate.
`, command, strings.Join(subcommands, " "))
	fs := NewFlagSet(usage)
	if err := fs.Parse(args); err != nil {
		return nil, fs
	}
	args = fs.Args()
	if len(args) == 0 {
		return nil, fs
	}

	switch args[0] {
	case "set":
		return CertificateSetArgsParser{}.Parse(command, append(subcommands, args[0]), args[1:])
	default:
		return nil, fs
	}
}

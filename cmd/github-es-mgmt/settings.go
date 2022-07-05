package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	mgmt "github.com/hnakamur/github-es-mgmt"
)

type SettingsArgsParser struct{}

func (p SettingsArgsParser) Parse(command string, subcommands, args []string) (Command, error) {
	usage := fmt.Sprintf(`Usage: %s %s <subcommand> [options]

subcommands:
    get       Get settings.
    set       Set settings.
`, command, strings.Join(subcommands, " "))
	fs := NewFlagSet(usage)
	if err := fs.Parse(args); err != nil {
		return nil, NewUsageError(fs, "")
	}
	args = fs.Args()
	if len(args) == 0 {
		return nil, NewUsageError(fs, "")
	}

	switch args[0] {
	case "get":
		return SettingsGetArgsParser{}.Parse(command, append(subcommands, args[0]), args[1:])
	case "set":
		return SettingsSetArgsParser{}.Parse(command, append(subcommands, args[0]), args[1:])
	default:
		return nil, NewUsageError(fs, "")
	}
}

func waitForConfigurationProcessToFinish(client *mgmt.Client, waitConfigInterval time.Duration) error {
	const configurationStatusSuccess = "success"

	for {
		s, err := client.GetConfigurationStatus()
		if err != nil {
			return err
		}
		sJson, err := json.Marshal(*s)
		if err != nil {
			return err
		}
		log.Printf("got configuration status: %s\n", sJson)
		if s.Status == configurationStatusSuccess {
			break
		}

		time.Sleep(waitConfigInterval)
	}
	return nil
}

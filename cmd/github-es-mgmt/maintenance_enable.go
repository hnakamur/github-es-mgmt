package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	mgmt "github.com/hnakamur/github-es-mgmt"
)

type MaintenanceEnableArgsParser struct{}

func (p MaintenanceEnableArgsParser) Parse(command string, subcommands, args []string) (Command, *flag.FlagSet, error) {
	usageTemplate := fmt.Sprintf(`Usage: %s %s <subcommand> [options]

options:
`, command, strings.Join(subcommands, " "))
	fs := newFlagSet(subcommands, usageTemplate)
	c := MaintenanceEnableCommand{}
	fs.StringVar(&c.Endpoint, "endpoint", "", "management API endpoint (ex. https://github-es.example.jp:8443)")
	fs.StringVar(&c.When, "when", "", "\"now\" or any date parsable by https://github.com/mojombo/chronic")
	fs.DurationVar(&c.Timeout, "timeout", 30*time.Second, "HTTP client timeout")
	if err := fs.Parse(args); err != nil {
		return nil, fs, nil
	}

	c.password = os.Getenv("MGMT_PASSWORD")
	if c.password == "" {
		return nil, fs, errors.New("Please set MGMT_PASSWORD environment variable")
	}
	if c.Endpoint == "" {
		return nil, fs, errors.New("Please set \"-endpoint\" flag")
	}
	if c.When == "" {
		return nil, fs, errors.New("Please set \"-when\" flag")
	}

	return &c, nil, nil
}

type MaintenanceEnableCommand struct {
	password string
	Endpoint string
	When     string
	Timeout  time.Duration
}

func (c *MaintenanceEnableCommand) Execute() error {
	cfg := mgmt.NewClientConfig().SetHTTPClient(&http.Client{Timeout: c.Timeout})
	client, err := mgmt.NewClient(c.Endpoint, c.password, cfg)
	if err != nil {
		return err
	}
	if err := client.EnableOrDisableMaintenanceMode(true, c.When); err != nil {
		return err
	}
	log.Printf("enabled maintenance mode successfully.")
	return nil
}

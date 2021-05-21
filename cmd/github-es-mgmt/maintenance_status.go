package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	mgmt "github.com/hnakamur/github-es-mgmt"
)

type MaintenanceStatusArgsParser struct{}

func (p MaintenanceStatusArgsParser) Parse(command string, subcommands, args []string) (Command, Usager) {
	usage := fmt.Sprintf(`Usage: %s %s [options]

options:
`, command, strings.Join(subcommands, " "))
	fs := NewFlagSet(usage)
	c := MaintenanceStatusCommand{}
	fs.StringVar(&c.Endpoint, "endpoint", "", "management API endpoint (ex. https://github-es.example.jp:8443)")
	fs.DurationVar(&c.Timeout, "timeout", 30*time.Second, "HTTP client timeout")
	if err := fs.Parse(args); err != nil {
		return nil, fs
	}

	c.password = os.Getenv("MGMT_PASSWORD")
	if c.password == "" {
		return nil, fs.SetError("Please set MGMT_PASSWORD environment variable")
	}

	if c.Endpoint == "" {
		return nil, fs.SetError("Please set \"-endpoint\" flag")
	}

	return &c, nil
}

type MaintenanceStatusCommand struct {
	password string
	Endpoint string
	Timeout  time.Duration
}

func (c *MaintenanceStatusCommand) Execute() error {
	cfg := mgmt.NewClientConfig().SetHTTPClient(&http.Client{Timeout: c.Timeout})
	client, err := mgmt.NewClient(c.Endpoint, c.password, cfg)
	if err != nil {
		return err
	}
	s, err := client.GetMaintenanceStatus()
	if err != nil {
		return err
	}
	log.Printf("got maintenance status: %+v", *s)
	return nil
}

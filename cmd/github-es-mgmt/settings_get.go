package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	mgmt "github.com/hnakamur/github-es-mgmt"
)

type SettingsGetArgsParser struct{}

func (p SettingsGetArgsParser) Parse(command string, subcommands, args []string) (Command, Usager) {
	usage := fmt.Sprintf(`Usage: %s %s <subcommand> [options]

options:
`, command, strings.Join(subcommands, " "))
	fs := NewFlagSet(usage)
	c := SettingsGetCommand{}
	fs.StringVar(&c.Endpoint, "endpoint", "", "management API endpoint (ex. https://github-es.example.jp:8443)")
	fs.StringVar(&c.Out, "out", "-", `output filename ("-" for stdout)`)
	fs.DurationVar(&c.Timeout, "timeout", 10*time.Minute, "HTTP client timeout")
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
	if c.Out == "" {
		return nil, fs.SetError("Please set \"-out\" flag")
	}

	return &c, nil
}

type SettingsGetCommand struct {
	password string
	Endpoint string
	Out      string
	Timeout  time.Duration
}

func (c *SettingsGetCommand) Execute() error {
	cfg := mgmt.NewClientConfig().SetHTTPClient(&http.Client{Timeout: c.Timeout})
	client, err := mgmt.NewClient(c.Endpoint, c.password, cfg)
	if err != nil {
		return err
	}
	settings, err := client.GetSettings()
	if err != nil {
		return err
	}
	if c.Out == "-" {
		fmt.Print(settings)
		return nil
	}
	return os.WriteFile(c.Out, []byte(settings), 0o644)
}

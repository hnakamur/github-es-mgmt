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

func (p SettingsGetArgsParser) Parse(command string, subcommands, args []string) (Command, error) {
	usage := fmt.Sprintf(`Usage: %s %s <subcommand> [options]

options:
`, command, strings.Join(subcommands, " "))
	fs := NewFlagSet(usage)
	c := SettingsGetCommand{}
	fs.StringVar(&c.Endpoint, "endpoint", "", "management API endpoint (ex. https://github-es.example.jp:8443)")
	fs.StringVar(&c.Out, "out", "-", `output filename ("-" for stdout)`)
	fs.DurationVar(&c.Timeout, "timeout", 10*time.Minute, "HTTP client timeout")
	fs.BoolVar(&c.TLSInsecureSkipVerify, "tls-insecure-skip-verify", false, "skip verify server's certificate. Use this only when server's certificate is expired.")
	if err := fs.Parse(args); err != nil {
		return nil, NewUsageError(fs, "")
	}

	c.password = GetManagementConsolePassword()

	if c.Endpoint == "" {
		return nil, NewUsageError(fs, "Please set \"-endpoint\" flag")
	}
	if c.Out == "" {
		return nil, NewUsageError(fs, "Please set \"-out\" flag")
	}

	return &c, nil
}

type SettingsGetCommand struct {
	password              string
	Endpoint              string
	Out                   string
	Timeout               time.Duration
	TLSInsecureSkipVerify bool
}

func (c *SettingsGetCommand) Execute() error {
	var roundTripper http.RoundTripper
	if c.TLSInsecureSkipVerify {
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.TLSClientConfig.InsecureSkipVerify = true
		roundTripper = transport
	}
	httpClient := &http.Client{Timeout: c.Timeout, Transport: roundTripper}

	cfg := mgmt.NewClientConfig().SetHTTPClient(httpClient)
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

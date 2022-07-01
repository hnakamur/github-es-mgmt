package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	mgmt "github.com/hnakamur/github-es-mgmt"
)

type SettingsSetArgsParser struct{}

func (p SettingsSetArgsParser) Parse(command string, subcommands, args []string) (Command, Usager) {
	usage := fmt.Sprintf(`Usage: %s %s <subcommand> [options]

options:
`, command, strings.Join(subcommands, " "))
	fs := NewFlagSet(usage)
	c := SettingsSetCommand{}
	fs.StringVar(&c.Endpoint, "endpoint", "", "management API endpoint (ex. https://github-es.example.jp:8443)")
	fs.StringVar(&c.In, "in", "-", `input filename ("-" for stdin)`)
	fs.DurationVar(&c.Timeout, "timeout", 30*time.Second, "HTTP client timeout")
	fs.DurationVar(&c.WaitConfigInterval, "interval", time.Minute, "polling interval for waiting configuration process to be finished")
	if err := fs.Parse(args); err != nil {
		return nil, fs
	}

	c.password = GetManagementConsolePassword()

	if c.Endpoint == "" {
		return nil, fs.SetError("Please set \"-endpoint\" flag")
	}
	if c.In == "" {
		return nil, fs.SetError("Please set \"-in\" flag")
	}
	return &c, nil
}

type SettingsSetCommand struct {
	password           string
	Endpoint           string
	In                 string
	Timeout            time.Duration
	WaitConfigInterval time.Duration
}

func (c *SettingsSetCommand) Execute() error {
	settings, err := readSettings(c.In)
	if err != nil {
		return err
	}

	cfg := mgmt.NewClientConfig().SetHTTPClient(&http.Client{Timeout: c.Timeout})
	client, err := mgmt.NewClient(c.Endpoint, c.password, cfg)
	if err != nil {
		return err
	}

	if err := client.SetSettings(string(settings)); err != nil {
		return fmt.Errorf("set settings: %s", err)
	}
	log.Printf("finished set settings API successfully.")

	if err := client.StartConfigurationProcess(); err != nil {
		return err
	}
	log.Printf("finished start configuration process API successfully.")

	if err := waitForConfigurationProcessToFinish(client, c.WaitConfigInterval); err != nil {
		return err
	}

	log.Printf("finished configuration process successfully.")
	return nil
}

func readSettings(in string) (settings []byte, err error) {
	if in == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(in)
}

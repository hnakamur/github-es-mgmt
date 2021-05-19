package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	mgmt "github.com/hnakamur/github-es-mgmt"
)

type SetMaintenanceCommand struct {
	password string
	Endpoint string
	Enabled  bool
	When     string
	Timeout  time.Duration
}

func (c *SetMaintenanceCommand) UsageTemplate() string {
	return `Usage: {{command}} set-maintenance [options]

options:
`
}

func (c *SetMaintenanceCommand) Parse(fs *flag.FlagSet, args []string) error {
	fs.StringVar(&c.Endpoint, "endpoint", "", "management API endpoint (ex. https://github-es.example.jp:8443)")
	fs.BoolVar(&c.Enabled, "enabled", true, "enable or disable maintenance mode")
	fs.StringVar(&c.When, "when", "", "\"now\" or any date parsable by https://github.com/mojombo/chronic")
	fs.DurationVar(&c.Timeout, "timeout", 10*time.Minute, "HTTP client timeout")
	if err := fs.Parse(args); err != nil {
		return err
	}

	c.password = os.Getenv("MGMT_PASSWORD")
	if c.password == "" {
		return errors.New("Please set MGMT_PASSWORD environment variable")
	}

	if c.Endpoint == "" {
		return errors.New("Please set \"-endpoint\" flag")
	}
	if _, err := url.Parse(c.Endpoint); err != nil {
		return fmt.Errorf("cannot parse endpoint: %s", err)
	}

	if c.When == "" {
		return errors.New("Please set \"-when\" flag")
	}

	return nil
}

func (c *SetMaintenanceCommand) Execute() error {
	cfg := mgmt.NewClientConfig().SetHTTPClient(&http.Client{Timeout: c.Timeout})
	client, err := mgmt.NewClient(c.Endpoint, c.password, cfg)
	if err != nil {
		return err
	}
	if err := client.EnableOrDisableMaintenanceMode(c.Enabled, c.When); err != nil {
		return err
	}
	var op string
	if c.Enabled {
		op = "enabled"
	} else {
		op = "disabled"
	}
	log.Printf("%s maintenance mode successfully.", op)
	return nil
}

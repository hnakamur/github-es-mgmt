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

type MaintenanceStatusCommand struct {
	password string
	Endpoint string
	Timeout  time.Duration
}

func (c *MaintenanceStatusCommand) UsageTemplate() string {
	return `Usage: {{command}} maintenance status [options]

options:
`
}

func (c *MaintenanceStatusCommand) Parse(fs *flag.FlagSet, args []string) error {
	fs.StringVar(&c.Endpoint, "endpoint", "", "management API endpoint (ex. https://github-es.example.jp:8443)")
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

	return nil
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

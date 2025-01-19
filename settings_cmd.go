package main

import (
	"fmt"
	"net/http"
)

type SettingsCmd struct {
	Get SettingsGetCmd `cmd:"" help:"Get the GHES Settings."`
}

type SettingsGetCmd struct {
	Endpoint string `required:"" env:"GHES_MANAGE_APIE_NDPOINT" help:"GHES Management API Endpoint (ex. https://your-github.example.jp:8443/manage)"`
	User     string `help:"A Management Console user name or \"api_key\" for the Root Site Administrator. Read from standard input when --user is not set"`
	Password string `help:"A password for a Management Console user or the Root Site Administrator. Read from standard input when --password is not set."`
}

func (c *SettingsGetCmd) Run(ctx *Context) error {
	user, password, err := getOrReadUserAndPassword(c.User, c.Password)
	if err != nil {
		return err
	}
	cli, err := newAPIClient(&http.Client{}, c.Endpoint, user, password)
	if err != nil {
		return err
	}

	respBody, err := cli.getSettings(ctx)
	if err != nil {
		return err
	}

	fmt.Print(string(respBody))
	return nil
}

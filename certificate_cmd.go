package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

type CertificateCmd struct {
	Get CertificateGetCmd `cmd:"" help:"Get certificate and key and save them to files"`
	Set CertificateSetCmd `cmd:"" help:"Set certificate and key"`
}

type CertificateSetCmd struct {
	Endpoint string `required:"" env:"GHES_MANAGE_APIE_NDPOINT" help:"GHES Management API Endpoint (ex. https://your-github.example.jp:8443/manage)"`
	User     string `help:"A Management Console user name or \"api_key\" for the Root Site Administrator. Read from standard input when --user is not set"`
	Password string `help:"A password for a Management Console user or the Root Site Administrator. Read from standard input when --password is not set."`

	Cert                     string        `required:"" help:"The certificate file"`
	Key                      string        `required:"" help:"The key file"`
	Apply                    bool          `help:"whether or not apply config change after setting certificate and key"`
	ApplyWaitPollingInterval time.Duration `default:"1m" help:"The polling interval for waiting config apply finished"`
}

func (c *CertificateSetCmd) Run(ctx *Context) error {
	user, password, err := getOrReadUserAndPassword(c.User, c.Password)
	if err != nil {
		return err
	}
	cli, err := newAPIClient(&http.Client{}, c.Endpoint, user, password)
	if err != nil {
		return err
	}

	cert, err := readFile(c.Cert)
	if err != nil {
		return err
	}
	key, err := readFile(c.Key)
	if err != nil {
		return err
	}
	if err := cli.setCertAndKey(ctx, cert, key); err != nil {
		return err
	}

	if !c.Apply {
		log.Printf("Exiting without applying the change of certificate and key. You need to apply the change by executing \"ghe-config-apply\" on your GitHub Enterprise Server.")
		return nil
	}

	log.Print("Applying the change of certificate and key...")
	startTime := time.Now()
	runID, err := cli.triggerConfigApply(ctx)
	if err != nil {
		return err
	}
	for {
		status, err := cli.getConfigApplyStatus(ctx, runID)
		if err != nil {
			return err
		}
		if !status.Running {
			if !status.Successful {
				return errors.New("failed to apply change of certificate and key")
			}
			log.Printf("finished setting certificate and key, elapsed=%s", time.Since(startTime))
			return nil
		}

		statusBytes, err := json.Marshal(status)
		if err != nil {
			return err
		}
		log.Printf("Waiting for the change to be applied, status:%s, will check again in %s...", string(statusBytes), c.ApplyWaitPollingInterval)

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(c.ApplyWaitPollingInterval):
		}
	}
}

func readFile(filename string) (content string, err error) {
	contentBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(contentBytes), nil
}

type CertificateGetCmd struct {
	Endpoint string `required:"" env:"GHES_MANAGE_APIE_NDPOINT" help:"GHES Management API Endpoint (ex. https://your-github.example.jp:8443/manage)"`
	User     string `help:"A Management Console user name or \"api_key\" for the Root Site Administrator. Read from standard input when --user is not set"`
	Password string `help:"A password for a Management Console user or the Root Site Administrator. Read from standard input when --password is not set."`

	Cert string `required:"" help:"The certificate file to save"`
	Key  string `required:"" help:"The key file to save"`
}

func (c *CertificateGetCmd) Run(ctx *Context) error {
	user, password, err := getOrReadUserAndPassword(c.User, c.Password)
	if err != nil {
		return err
	}
	cli, err := newAPIClient(&http.Client{}, c.Endpoint, user, password)
	if err != nil {
		return err
	}

	cert, key, err := cli.getCertAndKey(ctx)
	if err != nil {
		return err
	}

	if err := os.WriteFile(c.Cert, []byte(cert), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(c.Key, []byte(key), 0o400); err != nil {
		return err
	}
	return nil
}

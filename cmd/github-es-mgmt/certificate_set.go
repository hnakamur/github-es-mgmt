package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	mgmt "github.com/hnakamur/github-es-mgmt"
)

type CertificateSetArgsParser struct{}

func (p CertificateSetArgsParser) Parse(command string, subcommands, args []string) (Command, *flag.FlagSet, error) {
	usage := fmt.Sprintf(`Usage: %s %s <subcommand> [options]

options:
`, command, strings.Join(subcommands, " "))
	fs := newFlagSet(subcommands, usage)
	c := CertificateSetCommand{}
	fs.StringVar(&c.Endpoint, "endpoint", "", "management API endpoint (ex. https://github-es.example.jp:8443)")
	fs.StringVar(&c.CertFilename, "cert", "", "certificate PEM filename")
	fs.StringVar(&c.KeyFilename, "key", "", "key PEM filename")
	fs.DurationVar(&c.Timeout, "timeout", 10*time.Minute, "HTTP client timeout")
	fs.DurationVar(&c.WaitConfigInterval, "interval", time.Minute, "polling interval for waiting configuration process to be finished")
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
	if c.CertFilename == "" {
		return nil, fs, errors.New("Please set \"-cert\" flag")
	}
	if c.KeyFilename == "" {
		return nil, fs, errors.New("Please set \"-key\" flag")
	}
	return &c, nil, nil
}

type CertificateSetCommand struct {
	password           string
	Endpoint           string
	CertFilename       string
	KeyFilename        string
	Timeout            time.Duration
	WaitConfigInterval time.Duration
}

func (c *CertificateSetCommand) Execute() error {
	cert, key, err := readCertAndKey(c.CertFilename, c.KeyFilename)
	if err != nil {
		return err
	}

	cfg := mgmt.NewClientConfig().SetHTTPClient(&http.Client{Timeout: c.Timeout})
	client, err := mgmt.NewClient(c.Endpoint, c.password, cfg)
	if err != nil {
		return err
	}
	if err := setCertificate(client, cert, key); err != nil {
		return err
	}
	log.Printf("finished set settings API successfully.")

	if err := client.StartConfigurationProcess(); err != nil {
		return err
	}
	log.Printf("finished start configuration process API successfully.")

	for {
		s, err := client.GetConfigurationStatus()
		if err != nil {
			return err
		}
		log.Printf("got configuration status: %+v", *s)
		if s.Status == configurationStatusSuccess {
			break
		}

		time.Sleep(c.WaitConfigInterval)
	}

	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return err
	}
	remoteCertificates, err := getRemoteCerticates(u.Host)
	if err != nil {
		return err
	}
	log.Printf("NotBefore=%s, NotAfter=%s for certificate at %s",
		remoteCertificates[0].NotBefore.Format(time.RFC3339),
		remoteCertificates[0].NotAfter.Format(time.RFC3339),
		c.Endpoint,
	)

	return nil
}

type Settings struct {
	Enterprise Enterprise `json:"enterprise"`
}

type Enterprise struct {
	GithubSsl GithubSsl `json:"github_ssl"`
}

type GithubSsl struct {
	Enabled bool   `json:"enabled"`
	Cert    string `json:"cert"`
	Key     string `json:"key"`
}

func setCertificate(c *mgmt.Client, cert, key []byte) error {
	s, err := json.Marshal(Settings{
		Enterprise: Enterprise{
			GithubSsl: GithubSsl{
				Enabled: true,
				Cert:    string(cert),
				Key:     string(key),
			},
		},
	})
	if err != nil {
		return err
	}

	if err := c.SetSettings(string(s)); err != nil {
		return err
	}
	return nil
}

const configurationStatusSuccess = "success"

func readCertAndKey(certFilename, keyFilename string) (cert, key []byte, err error) {
	cert, err = os.ReadFile(certFilename)
	if err != nil {
		return nil, nil, fmt.Errorf("read certificate file: %s", err)
	}

	key, err = os.ReadFile(keyFilename)
	if err != nil {
		return nil, nil, fmt.Errorf("read key file: %s", err)
	}

	return cert, key, nil
}

func getRemoteCerticates(hostPort string) ([]*x509.Certificate, error) {
	conn, err := tls.Dial("tcp", hostPort, &tls.Config{})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	s := conn.ConnectionState()
	if len(s.PeerCertificates) == 0 {
		return nil, fmt.Errorf("no peer certificate at %s", hostPort)
	}

	return s.PeerCertificates, nil
}

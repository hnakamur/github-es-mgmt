package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	mgmt "github.com/hnakamur/github-es-mgmt"
)

func main() {
	endpoint := flag.String("endpoint", "", "management API endpoint (ex. https://github-es.example.jp:8443)")
	certFilename := flag.String("cert", "", "certificate PEM filename")
	keyFilename := flag.String("key", "", "key PEM filename")
	timeout := flag.Duration("timeout", 5*time.Minute, "HTTP client timeout")
	waitConfigInterval := flag.Duration("interval", time.Minute, "polling interval for waiting configuration process to be finished")
	flag.Parse()

	password := os.Getenv("MGMT_PASSWORD")
	if password == "" {
		flag.Usage()
		log.Fatal("Please set MGMT_PASSWORD environment variable")
	}

	if *endpoint == "" {
		flag.Usage()
		log.Fatal("Please set \"-endpoint\" flag")
	}
	if _, err := url.Parse(*endpoint); err != nil {
		log.Fatalf("cannot parse endpoint: %s", err)
	}

	if *certFilename == "" {
		flag.Usage()
		log.Fatal("Please set \"-cert\" flag")
	}
	if *keyFilename == "" {
		flag.Usage()
		log.Fatal("Please set \"-key\" flag")
	}

	if err := run(password, *endpoint, *certFilename, *keyFilename, *timeout, *waitConfigInterval); err != nil {
		log.Fatal(err)
	}
}

func run(password, endpoint, certFilename, keyFilename string, timeout, waitConfigInterval time.Duration) error {
	cert, key, err := readCertAndKey(certFilename, keyFilename)
	if err != nil {
		return err
	}

	cfg := mgmt.NewClientConfig().SetHTTPClient(&http.Client{Timeout: timeout})
	c, err := mgmt.NewClient(endpoint, password, cfg)
	if err != nil {
		return err
	}
	if err := setCertificate(c, cert, key); err != nil {
		return err
	}
	log.Printf("finished set settings API successfully.")

	if err := c.StartConfigurationProcess(); err != nil {
		return err
	}
	log.Printf("finished start configuration process API successfully.")

	for {
		s, err := c.GetConfigurationStatus()
		if err != nil {
			return err
		}
		log.Printf("got configuration status: %+v", *s)
		if s.Status == configurationStatusSuccess {
			break
		}

		time.Sleep(waitConfigInterval)
	}

	u, err := url.Parse(endpoint)
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
		endpoint,
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

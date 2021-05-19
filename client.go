package mgmt

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ClientConfig struct {
	client *http.Client
}

type Client struct {
	client   *http.Client
	endpoint *url.URL
	password string
}

func NewClientConfig() *ClientConfig {
	return &ClientConfig{
		client: &http.Client{},
	}
}

func (c *ClientConfig) SetHTTPClient(client *http.Client) *ClientConfig {
	c.client = client
	return c
}

func NewClient(endpoint, password string, config *ClientConfig) (*Client, error) {
	if config == nil {
		config = NewClientConfig()
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	return &Client{
		client:   config.client,
		endpoint: u,
		password: password,
	}, nil
}

func (c *Client) SetSettings(settingsJson string) error {
	data := url.Values{"settings": []string{settingsJson}}.Encode()
	r := strings.NewReader(data)
	u := url.URL{
		Scheme: c.endpoint.Scheme,
		Host:   c.endpoint.Host,
		Path:   "/setup/api/settings",
	}
	req, err := http.NewRequest("PUT", u.String(), r)
	if err != nil {
		return err
	}
	req.SetBasicAuth("api_key", c.password)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(data)))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < http.StatusOK || http.StatusMultipleChoices <= resp.StatusCode {
		return NewAPIError(resp.StatusCode, resp.Header, respBodyData)
	}
	return nil
}

func (c *Client) StartConfigurationProcess() error {
	u := url.URL{
		Scheme: c.endpoint.Scheme,
		Host:   c.endpoint.Host,
		Path:   "/setup/api/configure",
	}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth("api_key", c.password)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < http.StatusOK || http.StatusMultipleChoices <= resp.StatusCode {
		return NewAPIError(resp.StatusCode, resp.Header, respBodyData)
	}
	return nil
}

type ConfigurationStatus struct {
	Status   string                        `json:"status"`
	Progress []ConfigurationStatusProgress `json:"progress"`
}

type ConfigurationStatusProgress struct {
	Status string `json:"status"`
	Key    string `json:"key"`
}

func (c *Client) GetConfigurationStatus() (*ConfigurationStatus, error) {
	u := url.URL{
		Scheme: c.endpoint.Scheme,
		Host:   c.endpoint.Host,
		Path:   "/setup/api/configcheck",
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("api_key", c.password)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || http.StatusMultipleChoices <= resp.StatusCode {
		return nil, NewAPIError(resp.StatusCode, resp.Header, respBodyData)
	}

	s := &ConfigurationStatus{}
	if err := json.Unmarshal(respBodyData, s); err != nil {
		return nil, err
	}
	return s, nil
}

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
	reqBody := url.Values{"settings": []string{settingsJson}}.Encode()
	req, err := c.newRequest("PUT", "/setup/api/settings", reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(reqBody)))

	if _, err := c.doRequest(req); err != nil {
		return err
	}
	return nil
}

func (c *Client) StartConfigurationProcess() error {
	req, err := c.newRequest("POST", "/setup/api/configure", "")
	if err != nil {
		return err
	}
	if _, err := c.doRequest(req); err != nil {
		return err
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
	req, err := c.newRequest("GET", "/setup/api/configcheck", "")
	if err != nil {
		return nil, err
	}

	respBody, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	s := &ConfigurationStatus{}
	if err := json.Unmarshal(respBody, s); err != nil {
		return nil, err
	}
	return s, nil
}

type MaintenanceStatus struct {
	Status              string              `json:"status"`
	ScheduledTime       string              `json:"scheduled_time"`
	ConnectionServices  []ConnectionService `json:"connection_services"`
	CanUnsetMaintenance bool                `json:"can_unset_maintenance"`
}

type ConnectionService struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}

func (c *Client) GetMaintenanceStatus() (*MaintenanceStatus, error) {
	req, err := c.newRequest("GET", "/setup/api/maintenance", "")
	if err != nil {
		return nil, err
	}

	respBody, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	s := &MaintenanceStatus{}
	if err := json.Unmarshal(respBody, s); err != nil {
		return nil, err
	}
	return s, nil
}

func (c *Client) EnableOrDisableMaintenanceMode(enabled bool, when string) error {
	m := &struct {
		Enabled bool   `json:"enabled"`
		When    string `json:"when"`
	}{
		Enabled: enabled,
		When:    when,
	}
	maintenanceBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	reqBody := url.Values{"maintenance": []string{string(maintenanceBytes)}}.Encode()
	req, err := c.newRequest("POST", "/setup/api/maintenance", reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(reqBody)))

	if _, err := c.doRequest(req); err != nil {
		return err
	}
	return nil
}

func (c *Client) newRequest(method, path, body string) (*http.Request, error) {
	u := url.URL{
		Scheme: c.endpoint.Scheme,
		Host:   c.endpoint.Host,
		Path:   path,
	}
	var r io.Reader
	if body != "" {
		r = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, u.String(), r)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("api_key", c.password)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	return req, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || http.StatusMultipleChoices <= resp.StatusCode {
		return nil, NewAPIError(resp.StatusCode, resp.Header, body)
	}
	return body, nil
}

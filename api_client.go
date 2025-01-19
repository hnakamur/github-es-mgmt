package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type APIClient struct {
	httpClient        *http.Client
	manageEndpointURL *url.URL
	user              string
	password          string
}

func newAPIClient(httpClient *http.Client, managementEndpoint, user, password string) (*APIClient, error) {
	manageEndpointURL, err := url.Parse(managementEndpoint)
	if err != nil {
		return nil, err
	}

	return &APIClient{
		httpClient:        httpClient,
		manageEndpointURL: manageEndpointURL,
		user:              user,
		password:          password,
	}, nil
}

type GHESSettings struct {
	GithubSsl GHESSettingsGithubSsl `json:"github_ssl"`
}

type GHESSettingsGithubSsl struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

func newCertificateSetRequestBody(cert, key string) GHESSettings {
	return GHESSettings{
		GithubSsl: GHESSettingsGithubSsl{
			Cert: cert,
			Key:  key,
		},
	}
}

func (c *APIClient) setCertAndKey(ctx *Context, cert, key string) error {
	reqBodyBytes, err := json.Marshal(newCertificateSetRequestBody(cert, key))
	if err != nil {
		return err
	}
	reqBody := bytes.NewReader(reqBodyBytes)

	method := http.MethodPut

	requestURL := c.urlForPath("/v1/config/settings").String()
	resp, err := c.sendRequest(ctx, method, requestURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to send request:%s %s, err:%s", method, requestURL, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d (%s), response body: %s, request:%s %s",
			resp.StatusCode, resp.Status, strings.TrimSpace(string(respBody)), method, requestURL)
	}

	return nil
}

func (c *APIClient) triggerConfigApply(ctx *Context) (runID string, err error) {
	method := http.MethodPost
	requestURL := c.urlForPath("/v1/config/apply").String()
	resp, err := c.sendRequest(ctx, method, requestURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to send request:%s %s, err:%s",
			method, requestURL, err)
	}

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d (%s), response body: %s, request:%s %s",
			resp.StatusCode, resp.Status, strings.TrimSpace(string(respBody)), method, requestURL)
	}

	var respObj struct {
		RunID string `json:"run_id"`
	}
	if err := json.Unmarshal(respBody, &respObj); err != nil {
		return "", fmt.Errorf("failed to parse response body: %s, request:%s %s, err:%s",
			strings.TrimSpace(string(respBody)), method, requestURL, err)
	}

	return respObj.RunID, nil
}

type configApplyStatus struct {
	Running    bool                    `json:"running"`
	Successful bool                    `json:"successful"`
	Nodes      []configApplyStatusNode `json:"nodes"`
}

type configApplyStatusNode struct {
	RunID      string `json:"run_id"`
	Hostname   string `json:"hostname"`
	Running    bool   `json:"running"`
	Successful bool   `json:"successful"`
}

func (c *APIClient) getConfigApplyStatus(ctx *Context, runID string) (*configApplyStatus, error) {
	method := http.MethodGet
	var requestURL string
	{
		u := c.urlForPath("/v1/config/apply")
		u.RawQuery = "runID=" + url.QueryEscape(runID)
		requestURL = u.String()
	}
	resp, err := c.sendRequest(ctx, method, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request:%s %s, err:%s",
			method, requestURL, err)
	}

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d (%s), response body: %s, request:%s %s",
			resp.StatusCode, resp.Status, strings.TrimSpace(string(respBody)), method, requestURL)
	}

	var status configApplyStatus
	if err := json.Unmarshal(respBody, &status); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %s, request:%s %s, err:%s",
			strings.TrimSpace(string(respBody)), method, requestURL, err)
	}

	return &status, nil
}

func (c *APIClient) getCertAndKey(ctx *Context) (cert, key string, err error) {
	respBody, err := c.getSettings(ctx)
	if err != nil {
		return "", "", err
	}

	var settings GHESSettings
	if err := json.Unmarshal(respBody, &settings); err != nil {
		return "", "", err
	}
	return settings.GithubSsl.Cert, settings.GithubSsl.Key, nil
}

func (c *APIClient) getSettings(ctx *Context) (respBody []byte, err error) {
	method := http.MethodGet
	requestURL := c.urlForPath("/v1/config/settings").String()
	resp, err := c.sendRequest(ctx, method, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request:%s %s, err:%s",
			method, requestURL, err)
	}

	respBody, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d (%s), response body: %s, request:%s %s",
			resp.StatusCode, resp.Status, strings.TrimSpace(string(respBody)), method, requestURL)
	}
	return respBody, nil
}

func (c *APIClient) urlForPath(path ...string) *url.URL {
	return c.manageEndpointURL.JoinPath(path...)
}

func (c *APIClient) sendRequest(ctx context.Context, method, url string, reqBody io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.user, c.password)

	var httpClient http.Client
	resp, err = httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	resp.Body = io.NopCloser(bytes.NewReader(respBodyBytes))
	return resp, nil
}

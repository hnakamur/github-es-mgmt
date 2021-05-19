package mgmt

import (
	"fmt"
	"net/http"
)

type APIError struct {
	statusCode int
	header     http.Header
	body       []byte
}

func NewAPIError(statusCode int, header http.Header, body []byte) *APIError {
	return &APIError{
		statusCode: statusCode,
		header:     header,
		body:       body,
	}
}

func (e *APIError) StatusCode() int { return e.statusCode }

func (e *APIError) Header() http.Header { return e.header }

func (e *APIError) Body() []byte { return e.body }

func (e *APIError) Error() string {
	return fmt.Sprintf("API error, status=%d, header=%+v, body=%s",
		e.statusCode, e.header, string(e.body))
}

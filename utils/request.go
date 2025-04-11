package utils

import (
	"encoding/json"
	"net/http"
	"bytes"
)

type netHelper struct {}

// RequestOption defines a function type for configuring HTTP requests
type RequestOption func(*http.Request)

// WithToken adds an authorization token to the request
func WithToken(token string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set("Authorization", token)
	}
}

func (n netHelper) GET(url string, options ...RequestOption) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	
	// Apply all options
	for _, option := range options {
		option(req)
	}
	
	client := &http.Client{}
	return client.Do(req)
}

func (n netHelper) POST(url string, body any, options ...RequestOption) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	
	// Set default content-type if body is provided
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	// Apply all options
	for _, option := range options {
		option(req)
	}
	
	client := &http.Client{}
	return client.Do(req)
}

var NetHelper netHelper
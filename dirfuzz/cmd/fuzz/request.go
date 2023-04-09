package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Request represents an HTTP request.
type Request struct {
	Method string
	URL    string
	Body   []byte
	Header http.Header
}

// Do sends the HTTP request and returns the response.
func (r *Request) Do() (*http.Response, error) {
	req, err := http.NewRequest(r.Method, r.URL, bytes.NewReader(r.Body))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}

	req.Header = r.Header
	client := http.Client{}
	return client.Do(req)
}

// Send sends the HTTP request and returns the response body.
func (r *Request) Send() ([]byte, error) {
	resp, err := r.Do()
	if err != nil {
		return nil, fmt.Errorf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %v", err)
	}

	return body, nil
}

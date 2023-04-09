package fuzz

import (
    "bufio"
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "strings"

    "github.com/PuerkitoBio/goquery"
)

// Scanner is responsible for generating fuzzing requests.
type Scanner struct {
    baseURL      string
    inputDir     string
    cookieHeader string
    filters      []*Filter
    client       *http.Client
    // debugFunc is a function that will be called to print debug messages.
    debugFunc func(msg string)
}

// NewScanner returns a new Scanner instance.
func NewScanner(baseURL string, inputDir string, cookieHeader string, filters []*Filter, client *http.Client, debugFunc func(msg string)) *Scanner {
    return &Scanner{
        baseURL:      baseURL,
        inputDir:     inputDir,
        cookieHeader: cookieHeader,
        filters:      filters,
        client:       client,
        debugFunc:    debugFunc,
    }
}

// Run executes the scanner.
func (s *Scanner) Run() error {
    err := filepath.Walk(s.inputDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.Mode().IsRegular() {
            return nil
        }
        if filepath.Ext(path) != ".txt" {
            return nil
        }

        // Open file
        f, err := os.Open(path)
        if err != nil {
            return err
        }
        defer f.Close()

        // Read file contents
        scanner := bufio.NewScanner(f)
        for scanner.Scan() {
            // Trim leading/trailing spaces
            line := strings.TrimSpace(scanner.Text())

            // Skip empty lines
            if line == "" {
                continue
            }

            // Skip comment lines
            if strings.HasPrefix(line, "#") {
                continue
            }

            // Make request
            err := s.makeRequest(line)
            if err != nil {
                s.debugFunc(fmt.Sprintf("[ERROR] %v", err))
            }
        }

        if err := scanner.Err(); err != nil {
            return err
        }

        return nil
    })

    return err
}

// makeRequest sends an HTTP request with the given payload.
func (s *Scanner) makeRequest(payload string) error {
	// Apply filters to payload
	for _, filter := range s.filters {
		if !filter.Match(payload) {
			return nil
		}
	}

	// Prepare request
	req, err := http.NewRequest(s.method, s.url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return err
	}

	// Set request headers
	for k, v := range s.headers {
		req.Header.Set(k, v)
	}

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Parse response body with goquery
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return err
	}

	// Print response body
	fmt.Printf("%s\n", doc.Text())

	return nil
}



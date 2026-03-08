// Package c64u provides a Go client for the Ultimate 1541-II+ REST API.
//
// The API is available starting from Ultimate firmware 3.11.
// See https://1541u-documentation.readthedocs.io/en/latest/api/api_calls.html
package c64u

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// Client communicates with the Ultimate 1541-II+ REST API.
type Client struct {
	// BaseURL is the base URL of the Ultimate device, e.g. "http://192.168.1.100".
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new API client for the Ultimate device at the given address.
// The address should be the IP or hostname, e.g. "192.168.1.100".
func NewClient(address string) *Client {
	return &Client{
		BaseURL:    "http://" + address,
		HTTPClient: http.DefaultClient,
	}
}

// ErrorResponse is the common error structure returned by the API.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

func (c *Client) url(path string, params url.Values) string {
	u := c.BaseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	return u
}

func (c *Client) doJSON(method, path string, params url.Values, body io.Reader, contentType string, result any) error {
	u := c.url(path, params)
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}
	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}
	return nil
}

func (c *Client) get(path string, params url.Values, result any) error {
	return c.doJSON(http.MethodGet, path, params, nil, "", result)
}

func (c *Client) put(path string, params url.Values, result any) error {
	return c.doJSON(http.MethodPut, path, params, nil, "", result)
}

func (c *Client) postJSON(path string, params url.Values, data any, result any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshaling JSON body: %w", err)
	}
	return c.doJSON(http.MethodPost, path, params, bytes.NewReader(body), "application/json", result)
}

func (c *Client) postFile(path string, params url.Values, filename string, fileData []byte, result any) error {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("creating form file: %w", err)
	}
	if _, err := part.Write(fileData); err != nil {
		return fmt.Errorf("writing file data: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("closing multipart writer: %w", err)
	}
	return c.doJSON(http.MethodPost, path, params, &buf, w.FormDataContentType(), result)
}

func (c *Client) getRaw(path string, params url.Values) ([]byte, error) {
	u := c.url(path, params)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}
	return io.ReadAll(resp.Body)
}

func (c *Client) postBinary(path string, params url.Values, data []byte, result any) error {
	return c.doJSON(http.MethodPost, path, params, bytes.NewReader(data), "application/octet-stream", result)
}

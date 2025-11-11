package comdirect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

// Client ist ein leichter Wrapper um http.Client mit nützlichen Helfern.
type Client struct {
	http *http.Client
}

// NewClient erstellt einen konfigurierten HTTP-Client.
// timeout: Gesamt-Timeout für Requests; nil cookiejar bedeutet ohne Cookies.
func NewClient(timeout time.Duration) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	c := &http.Client{
		Timeout: timeout,
		Jar:     jar,
		// Transport kann hier weiter konfiguriert werden (TLS, Proxies, KeepAlive, MaxIdleConnsPerHost usw.)
	}
	return &Client{http: c}, nil
}

// Do führt ein *http.Request aus und liefert den Body (schließt response.Body).
func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, []byte, error) {
	req = req.WithContext(ctx)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}
	return resp, body, nil
}

// GetJSON macht eine GET-Anfrage und unmarshalt das JSON in out.
func (c *Client) GetJSON(ctx context.Context, url string, out interface{}) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Accept", "application/json")
	resp, body, err := c.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return resp, fmt.Errorf("http %d: %s", resp.StatusCode, string(body))
	}
	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

// PostJSON sendet JSON und optional unmarshalt die Antwort in out.
func (c *Client) PostJSON(ctx context.Context, url string, payload interface{}, out interface{}) (*http.Response, error) {
	var buf bytes.Buffer
	if payload != nil {
		if err := json.NewEncoder(&buf).Encode(payload); err != nil {
			return nil, err
		}
	}
	req, _ := http.NewRequest(http.MethodPost, url, &buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, body, err := c.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return resp, fmt.Errorf("http %d: %s", resp.StatusCode, string(body))
	}
	if out != nil && len(body) > 0 {
		if err := json.Unmarshal(body, out); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

// PostForm sendet Formulardaten (application/x-www-form-urlencoded)
// und optional unmarshallt die Antwort in out.
// data can be provided as url.Values. If nil provided, an empty body is sent.
func (c *Client) PostForm(ctx context.Context, urlStr string, data url.Values, out interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if data != nil {
		encoded := data.Encode()
		bodyReader = strings.NewReader(encoded)
	} else {
		bodyReader = bytes.NewReader(nil)
	}

	req, _ := http.NewRequest(http.MethodPost, urlStr, bodyReader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, respBody, err := c.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return resp, fmt.Errorf("http %d: %s", resp.StatusCode, string(respBody))
	}
	if out != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, out); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

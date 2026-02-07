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
	// defaultHeader are headers applied to every request if not already set on the request.
	defaultHeader http.Header
}

// Context key and helpers for per-request headers stored in context.
type ctxHeadersKeyType struct{}

// WithHeaders returns a new context that carries a copy of the provided headers.
func WithHeaders(ctx context.Context, h http.Header) context.Context {
	if h == nil {
		return ctx
	}
	copy := make(http.Header, len(h))
	for k, vals := range h {
		copy[k] = append([]string(nil), vals...)
	}
	return context.WithValue(ctx, ctxHeadersKeyType{}, copy)
}

// headersFromCtx extracts headers put into context via WithHeaders.
func headersFromCtx(ctx context.Context) http.Header {
	if ctx == nil {
		return nil
	}
	if v := ctx.Value(ctxHeadersKeyType{}); v != nil {
		if h, ok := v.(http.Header); ok {
			return h
		}
	}
	return nil
}

// SetDefaultHeaders sets client-wide headers. The header map is copied.
func (c *Client) SetDefaultHeaders(h http.Header) {
	if c == nil {
		return
	}
	if h == nil {
		c.defaultHeader = nil
		return
	}
	copy := make(http.Header, len(h))
	for k, vals := range h {
		copy[k] = append([]string(nil), vals...)
	}
	c.defaultHeader = copy
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
	// 1) inject headers from context (per-request). These have lower priority than explicit req headers.
	if hdrs := headersFromCtx(req.Context()); hdrs != nil {
		for k, vals := range hdrs {
			if req.Header.Get(k) == "" {
				for _, v := range vals {
					req.Header.Add(k, v)
				}
			}
		}
	}

	// 2) apply client default headers for any header not already set on the request
	if c != nil && c.defaultHeader != nil {
		for k, vals := range c.defaultHeader {
			if req.Header.Get(k) == "" {
				for _, v := range vals {
					req.Header.Add(k, v)
				}
			}
		}
	}

	// Headers auf Console ausgeben (Debugging)
	fmt.Printf("Request: %s %s\n", req.Method, req.URL)
	for k, vals := range req.Header {
		for _, v := range vals {
			fmt.Printf("  %s: %s\n", k, v)
		}
	}
	fmt.Println()

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
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	req.Header.Set("Content-Type", "application/json")
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

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, urlStr, bodyReader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

// PatchJSON sendet JSON per PATCH und optional unmarshalt die Antwort in out.
func (c *Client) PatchJSON(ctx context.Context, url string, payload interface{}, out interface{}) (*http.Response, error) {
	var buf bytes.Buffer
	if payload != nil {
		if err := json.NewEncoder(&buf).Encode(payload); err != nil {
			return nil, err
		}
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodPatch, url, &buf)
	req.Header.Set("Content-Type", "application/json")
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

func (c *Client) Delete(ctx context.Context, url string) (*http.Response, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	resp, _, err := c.Do(ctx, req)
	return resp, err
}

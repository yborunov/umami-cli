package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type Client struct {
	baseURL    *url.URL
	token      string
	httpClient *http.Client
}

func New(endpoint, token string) (*Client, error) {
	if endpoint == "" {
		return nil, errors.New("endpoint required")
	}
	if !strings.Contains(endpoint, "://") {
		return nil, fmt.Errorf("endpoint must include scheme: %s", endpoint)
	}

	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint: %w", err)
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/")

	return &Client{
		baseURL: parsed,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (c *Client) WithToken(token string) *Client {
	clone := *c
	clone.token = token
	return &clone
}

func (c *Client) Do(ctx context.Context, method, p string, body any, out any, auth bool) (int, error) {
	url, err := c.buildURL(p)
	if err != nil {
		return 0, err
	}

	var reader io.Reader
	if body != nil {
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return 0, err
		}
		reader = buf
	}

	req, err := http.NewRequestWithContext(ctx, method, url.String(), reader)
	if err != nil {
		return 0, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if auth {
		if c.token == "" {
			return 0, errors.New("missing token: run `umami auth login` or set UMAMI_TOKEN")
		}
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if out != nil {
		req.Header.Set("Accept", "application/json")
	}

	if debugEnabled() {
		fmt.Fprintf(os.Stderr, "debug: http request method=%s url=%s auth=%t token-set=%t token-len=%d\n",
			method, url.String(), auth, c.token != "", len(c.token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if debugEnabled() {
		fmt.Fprintf(os.Stderr, "debug: http response status=%d url=%s\n", resp.StatusCode, url.String())
	}

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
		return resp.StatusCode, fmt.Errorf("request failed (%d): %s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}

	if out == nil {
		return resp.StatusCode, nil
	}
	return resp.StatusCode, json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) DoRaw(ctx context.Context, method, p string, body any, auth bool) (int, []byte, error) {
	url, err := c.buildURL(p)
	if err != nil {
		return 0, nil, err
	}

	var reader io.Reader
	if body != nil {
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return 0, nil, err
		}
		reader = buf
	}

	req, err := http.NewRequestWithContext(ctx, method, url.String(), reader)
	if err != nil {
		return 0, nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if auth {
		if c.token == "" {
			return 0, nil, errors.New("missing token: run `umami auth login` or set UMAMI_TOKEN")
		}
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "application/json")

	if debugEnabled() {
		fmt.Fprintf(os.Stderr, "debug: http request method=%s url=%s auth=%t token-set=%t token-len=%d\n",
			method, url.String(), auth, c.token != "", len(c.token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	if debugEnabled() {
		fmt.Fprintf(os.Stderr, "debug: http response status=%d url=%s\n", resp.StatusCode, url.String())
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}

	if resp.StatusCode >= 400 {
		return resp.StatusCode, bodyBytes, fmt.Errorf("request failed (%d): %s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}
	return resp.StatusCode, bodyBytes, nil
}

func (c *Client) buildURL(p string) (*url.URL, error) {
	url := *c.baseURL
	pathPart := p
	query := ""
	if idx := strings.Index(p, "?"); idx != -1 {
		pathPart = p[:idx]
		query = p[idx+1:]
	}
	if pathPart != "" {
		url.Path = path.Join(url.Path, pathPart)
	}
	if query != "" {
		url.RawQuery = query
	}
	return &url, nil
}

func debugEnabled() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("DEBUG")), "true")
}

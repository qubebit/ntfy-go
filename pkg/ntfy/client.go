package ntfy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type (
	Client struct {
		httpClient *http.Client
		host       *url.URL
		headers    http.Header
	}

	PublishOpts struct {
		Message *Message
	}

	PublishResult struct {
		ID      string `json:"id"`      // :"bUhbhgmmbeW0"
		Time    int    `json:"time"`    // :1685150791
		Expires int    `json:"expires"` // :1685193991
		Event   string `json:"event"`   // :"message"
		Topic   string `json:"topic"`   // :"Server"
		Message string `json:"message"` // :"triggered"
	}

	Options struct {
		HTTPClient *http.Client
		Headers    http.Header
		Host       string
	}

	Option func(*Options)
)

var (
	ErrNoServer = errors.New("server is nil")
	ErrNoTopic  = errors.New("topic is nil")
)

// New creates a ntfy client with the given options
func New(opts ...Option) (*Client, error) {
	options := &Options{
		HTTPClient: http.DefaultClient,
		Headers:    http.Header{"Content-Type": []string{"application/json"}},
		Host:       "https://ntfy.sh",
	}

	for _, o := range opts {
		o(options)
	}

	host, err := url.Parse(options.Host)
	if err != nil {
		return nil, err
	}

	return &Client{
		httpClient: options.HTTPClient,
		host:       host,
		headers:    options.Headers,
	}, nil
}

func (t *Client) Publish(ctx context.Context, opts *PublishOpts) (*PublishResult, error) {
	buf, err := json.Marshal(opts.Message)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.host.String(), bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if s := resp.StatusCode; s < 200 || s >= 300 {
		return nil, fmt.Errorf("non-200 http response code from server: %d", s)
	}

	var pubResp PublishResult
	if err = json.NewDecoder(resp.Body).Decode(&pubResp); err != nil {
		return nil, err
	}

	return &pubResp, nil
}

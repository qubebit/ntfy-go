package ntfy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-playground/validator/v10"
)

type (
	Client struct {
		httpClient *http.Client
		validator  *validator.Validate
		host       *url.URL
		headers    http.Header
	}

	PublishOpts struct {
		Message *Message    `validate:"required"`
		Headers http.Header `validate:"-"`
	}

	PublishResult struct {
		ID      string `json:"id"`
		Time    int    `json:"time"`
		Expires int    `json:"expires"`
		Event   string `json:"event"`
		Topic   string `json:"topic"`
		Message string `json:"message"`
	}

	Options struct {
		HTTPClient *http.Client
		Validator  *validator.Validate
		Headers    http.Header
		Host       string
	}

	Option func(*Options)
)

var (
	ErrMissingHost       = errors.New("missing host")
	ErrMissingHTTPClient = errors.New("missing http client")
	ErrMissingValidator  = errors.New("missing validator")
)

func WithHTTPClient(client *http.Client) Option {
	return func(o *Options) {
		o.HTTPClient = client
	}
}

func WithValidator(v *validator.Validate) Option {
	return func(o *Options) {
		o.Validator = v
	}
}

func WithHeaders(headers http.Header) Option {
	return func(o *Options) {
		o.Headers = headers
	}
}

func WithHost(host string) Option {
	return func(o *Options) {
		o.Host = host
	}
}

// New creates a ntfy client with the given options
func New(opts ...Option) (*Client, error) {
	options := &Options{
		HTTPClient: http.DefaultClient,
		Validator:  validator.New(),
		Host:       "https://ntfy.sh",
		Headers:    http.Header{"Content-Type": []string{"application/json"}},
	}
	for _, o := range opts {
		o(options)
	}

	if options.HTTPClient == nil {
		return nil, ErrMissingHTTPClient
	}

	if options.Validator == nil {
		return nil, ErrMissingValidator
	}

	if options.Host == "" {
		return nil, ErrMissingHost
	}

	host, err := url.Parse(options.Host)
	if err != nil {
		return nil, err
	}

	return &Client{
		httpClient: options.HTTPClient,
		validator:  options.Validator,
		headers:    options.Headers,
		host:       host,
	}, nil
}

func (c *Client) Publish(ctx context.Context, opts *PublishOpts) (*PublishResult, error) {
	if err := c.validator.Struct(opts); err != nil {
		return nil, err
	}

	buf, err := json.Marshal(opts.Message)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.host.String(), bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	for key, values := range c.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	if opts.Headers != nil {
		for key, values := range opts.Headers {
			req.Header.Del(key)
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

	resp, err := c.httpClient.Do(req)
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

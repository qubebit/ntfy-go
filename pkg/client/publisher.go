package client

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
	// Publisher creates messages for topics
	Publisher struct {
		server     *url.URL
		httpClient *http.Client

		Headers http.Header
	}

	PublishResult struct {
		ID      string `json:"id"`      // :"bUhbhgmmbeW0"
		Time    int    `json:"time"`    // :1685150791
		Expires int    `json:"expires"` // :1685193991
		Event   string `json:"event"`   // :"message"
		Topic   string `json:"topic"`   // :"Server"
		Message string `json:"message"` // :"triggered"
	}
)

var (
	ErrNoServer = errors.New("server is nil")
	ErrNoTopic  = errors.New("topic is nil")
)

// NewPublisher creates a topic publisher for the specified server URL,
// and uses the supplied HTTP client to resolve the request
func NewPublisher(server *url.URL, httpClient *http.Client) (*Publisher, error) {
	if server == nil {
		return nil, ErrNoServer
	}

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Publisher{
		server:     server,
		httpClient: httpClient,
		Headers:    http.Header{"Content-Type": []string{"application/json"}},
	}, nil
}

func (t *Publisher) SendMessage(ctx context.Context, m *Message) (*PublishResult, error) {
	buf, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.server.String(), bytes.NewReader(buf))
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

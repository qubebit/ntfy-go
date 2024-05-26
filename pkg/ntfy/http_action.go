package ntfy

import (
	"encoding/json"
	"net/url"
)

type (
	// HttpAction allows attaching an HTTP request action to a notification
	HttpAction[X comparable] struct {
		Label   string            // Label of the action button in the notification
		URL     *url.URL          // URL to which the HTTP request will be sent
		Method  string            // HTTP method to use for request, default is POST
		Headers map[string]string // HTTP headers to pass in request
		Body    X                 // HTTP body
		Clear   bool              // Clear notification after HTTP request succeeds
	}

	httpAction struct {
		Action  string            `json:"action"`
		Label   string            `json:"label"`
		URL     string            `json:"url"`
		Method  string            `json:"method,omitempty"`
		Headers map[string]string `json:"headers,omitempty"`
		Body    string            `json:"body,omitempty"`
		Clear   bool              `json:"clear,omitempty"`
	}
)

func (h *HttpAction[X]) actionType() ActionButtonType {
	return HTTP
}

func (h *HttpAction[X]) MarshalJSON() ([]byte, error) {
	url := ""
	if h.URL != nil {
		url = h.URL.String()
	}

	body := ""

	var zeroVal X
	if h.Body != zeroVal {
		b, err := json.Marshal(h.Body)
		if err != nil {
			return nil, err
		}

		body = string(b)
	}

	return json.Marshal(&httpAction{
		Action:  "http",
		Label:   h.Label,
		URL:     url,
		Method:  h.Method,
		Headers: h.Headers,
		Clear:   h.Clear,
		Body:    body,
	})
}

package ntfy

import (
	"encoding/json"
	"net/url"
)

type (
	ViewAction struct {
		Label string
		Link  *url.URL
		Clear bool
	}

	viewAction struct {
		Action string `json:"action"`
		Label  string `json:"label"`
		URL    string `json:"url,omitempty"`
		Clear  bool   `json:"clear,omitempty"`
	}
)

func (v *ViewAction) actionType() ActionButtonType {
	return View
}

func (v *ViewAction) MarshalJSON() ([]byte, error) {
	url := ""
	if v.Link != nil {
		url = v.Link.String()
	}

	return json.Marshal(&viewAction{
		Action: "view",
		Label:  v.Label,
		URL:    url,
		Clear:  v.Clear,
	})
}

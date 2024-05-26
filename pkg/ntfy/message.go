package ntfy

import (
	"encoding/json"
	"net/url"
	"time"
)

type (
	// Message is a struct you can create from TopicPublisher that
	// will publish a message to the specified topic. This method does not allow
	// for attaching files to the notification, but it can post a link to an attachment
	Message struct {
		Topic    string         `validate:"required"` // Target topic name
		Message  string         // Message body; set to triggered if empty or not passed
		Title    string         // Message title
		Tags     []string       // List of tags that may or not map to emojis
		Priority Priority       // Message priority with 1=min, 3=default and 5=max
		Actions  []ActionButton // Custom user action buttons for notifications
		ClickURL *url.URL       // Website opened when notification is clicked
		IconURL  *url.URL       // URL to use as notification icon
		Delay    time.Duration  // Duration to delay delivery
		Email    string         // E-mail address for e-mail notifications
		Call     string         // Phone number to use for voice call

		AttachURLFilename string   // File name of the attachment
		AttachURL         *url.URL // URL of an attachment
	}

	message struct {
		Topic     string         `json:"topic"`
		Message   string         `json:"message,omitempty"`
		Title     string         `json:"title,omitempty"`
		Tags      []string       `json:"tags,omitempty"`
		Priority  Priority       `json:"priority,omitempty"`
		Actions   []ActionButton `json:"actions,omitempty"`
		Click     string         `json:"click,omitempty"`
		Icon      string         `json:"icon,omitempty"`
		Delay     string         `json:"delay,omitempty"`
		Email     string         `json:"email,omitempty"`
		Call      string         `json:"call,omitempty"`
		Filename  string         `json:"filename,omitempty"`
		AttachURL string         `json:"attachurl,omitempty"`
	}
)

func (m *Message) MarshalJSON() ([]byte, error) {
	click := ""
	if m.ClickURL != nil {
		click = m.ClickURL.String()
	}

	icon := ""
	if m.IconURL != nil {
		icon = m.IconURL.String()
	}

	attachURL := ""
	if m.AttachURL != nil {
		attachURL = m.AttachURL.String()
	}

	delay := ""
	if m.Delay > 0 {
		delay = m.Delay.String()
	}

	priority := m.Priority
	if priority < 0 {
		priority = 0
	}

	return json.Marshal(message{
		Topic:     m.Topic,
		Message:   m.Message,
		Title:     m.Title,
		Tags:      m.Tags,
		Priority:  priority,
		Actions:   m.Actions,
		Click:     click,
		Icon:      icon,
		Delay:     delay,
		Email:     m.Email,
		Call:      m.Call,
		Filename:  m.AttachURLFilename,
		AttachURL: attachURL,
	})
}

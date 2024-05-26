package ntfy

import (
	"encoding/json"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestMessageMarshalJSON(t *testing.T) {
	testCases := []struct {
		name        string
		arg         Message
		expected    message
		expectedErr error
	}{
		{
			name:     "Empty Topic",
			expected: message{Topic: ""},
		},
		{
			name:     "Non-empty Topic",
			arg:      Message{Topic: "topic"},
			expected: message{Topic: "topic"},
		},
		{
			name:     "Message Field",
			arg:      Message{Message: "Message"},
			expected: message{Topic: "", Message: "Message"},
		},
		{
			name:     "Title Field",
			arg:      Message{Title: "Title"},
			expected: message{Topic: "", Title: "Title"},
		},
		{
			name:     "Tags Field",
			arg:      Message{Tags: []string{"tag1", "tag2"}},
			expected: message{Topic: "", Tags: []string{"tag1", "tag2"}},
		},
		{
			name:     "Negative Priority",
			arg:      Message{Priority: -1},
			expected: message{Topic: "", Priority: 0},
		},
		{
			name:     "Positive Priority",
			arg:      Message{Priority: 1},
			expected: message{Topic: "", Priority: 1},
		},
		{
			name: "Action Buttons",
			arg: Message{Actions: []ActionButton{&ViewAction{
				Label: "view",
				Link:  &url.URL{Scheme: "http", Host: "example.com"},
				Clear: true,
			}}},
			expected: message{Topic: "", Actions: []ActionButton{&ViewAction{
				Label: "view",
				Link:  &url.URL{Scheme: "http", Host: "example.com"},
				Clear: true,
			}}},
		},
		{
			name:     "Click URL",
			arg:      Message{ClickURL: &url.URL{Scheme: "https", Host: "example.com"}},
			expected: message{Topic: "", Click: "https://example.com"},
		},
		{
			name:     "Icon URL",
			arg:      Message{IconURL: &url.URL{Scheme: "https", Host: "example.com"}},
			expected: message{Topic: "", Icon: "https://example.com"},
		},
		{
			name:     "Delay Field",
			arg:      Message{Delay: 1},
			expected: message{Topic: "", Delay: "1ns"},
		},
		{
			name:     "Email Field",
			arg:      Message{Email: "email@example.com"},
			expected: message{Topic: "", Email: "email@example.com"},
		},
		{
			name:     "Call Field",
			arg:      Message{Call: "1234567890"},
			expected: message{Topic: "", Call: "1234567890"},
		},
		{
			name:     "Attachment URL Filename",
			arg:      Message{AttachURLFilename: "file.txt"},
			expected: message{Topic: "", Filename: "file.txt"},
		},
		{
			name:     "Attachment URL",
			arg:      Message{AttachURL: &url.URL{Scheme: "https", Host: "example.com"}},
			expected: message{Topic: "", AttachURL: "https://example.com"},
		},
		{
			name: "All Fields",
			arg: Message{
				Topic:    "Topic",
				Message:  "Message",
				Title:    "Title",
				Tags:     []string{"tag1", "tag2"},
				Priority: 4,
				Actions: []ActionButton{&ViewAction{
					Label: "view",
					Link:  &url.URL{Scheme: "https", Host: "example.com"},
					Clear: true,
				}},
				ClickURL:          &url.URL{Scheme: "https", Host: "example.com"},
				IconURL:           &url.URL{Scheme: "https", Host: "example.com"},
				Delay:             100 * time.Nanosecond,
				Email:             "email@example.com",
				Call:              "1234567890",
				AttachURLFilename: "file.txt",
				AttachURL:         &url.URL{Scheme: "https", Host: "example.com"},
			},
			expected: message{
				Topic:    "Topic",
				Message:  "Message",
				Title:    "Title",
				Tags:     []string{"tag1", "tag2"},
				Priority: 4,
				Actions: []ActionButton{&ViewAction{
					Label: "view",
					Link:  &url.URL{Scheme: "https", Host: "example.com"},
					Clear: true,
				}},
				Click:     "https://example.com",
				Icon:      "https://example.com",
				Delay:     "100ns",
				Email:     "email@example.com",
				Call:      "1234567890",
				Filename:  "file.txt",
				AttachURL: "https://example.com",
			},
		},
		{
			name: "Specific Test Case",
			arg: Message{
				Topic:    "Sample topic",
				Tags:     []string{"sample"},
				Priority: 0,
				ClickURL: &url.URL{Scheme: "https", Host: "example.com"},
				IconURL:  &url.URL{Scheme: "https", Host: "example.com"},
				Message:  "Sample message content",
			},
			expected: message{
				Topic:   "Sample topic",
				Tags:    []string{"sample"},
				Message: "Sample message content",
				Click:   "https://example.com",
				Icon:    "https://example.com",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotBytes, err := json.Marshal(&tc.arg)
			if (err != nil) != (tc.expectedErr != nil) {
				t.Fatalf("unexpected error state: got error = %v, expected error = %v", err, tc.expectedErr)
			}

			if err != nil {
				if err.Error() != tc.expectedErr.Error() {
					t.Fatalf("unexpected error message: got error = %v, expected error = %v", err, tc.expectedErr)
				}
				return
			}

			expectedBytes, err := json.Marshal(tc.expected)
			if err != nil {
				t.Fatalf("failed to marshal expected value: %v", err)
			}

			if !reflect.DeepEqual(gotBytes, expectedBytes) {
				t.Errorf("unexpected result:\n got  = %s\n want = %s", string(gotBytes), string(expectedBytes))
			}
		})
	}
}

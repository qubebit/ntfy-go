package ntfy

import (
	"encoding/json"
	"net/url"
	"reflect"
	"testing"
)

func TestHTTPActionMarshal(t *testing.T) {
	tests := []struct {
		name    string
		action  HttpAction[any]
		want    httpAction
		wantErr bool
	}{
		{
			name: "Simple case",
			action: HttpAction[any]{
				Label:  "Test",
				URL:    &url.URL{Scheme: "https", Host: "example.com", Path: "/test"},
				Method: "GET",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body:  map[string]string{"key": "value"},
				Clear: true,
			},
			want: httpAction{
				Action: "http",
				Label:  "Test",
				URL:    "https://example.com/test",
				Method: "GET",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body:  "{\"key\":\"value\"}",
				Clear: true,
			},
			wantErr: false,
		},
		{
			name: "Empty URL and POST method",
			action: HttpAction[any]{
				Label:  "Empty URL",
				Method: "POST",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body:  map[string]string{"key": "value"},
				Clear: false,
			},
			want: httpAction{
				Action: "http",
				Label:  "Empty URL",
				Method: "POST",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body:  "{\"key\":\"value\"}",
				Clear: false,
			},
			wantErr: false,
		},
		{
			name: "Empty body",
			action: HttpAction[any]{
				Label:  "No Body",
				URL:    &url.URL{Scheme: "https", Host: "example.com", Path: "/test"},
				Method: "PUT",
				Clear:  true,
			},
			want: httpAction{
				Action: "http",
				Label:  "No Body",
				URL:    "https://example.com/test",
				Method: "PUT",
				Clear:  true,
			},
			wantErr: false,
		},
		{
			name: "Nil URL",
			action: HttpAction[any]{
				Label:  "Nil URL",
				Method: "DELETE",
				Body:   "simple string body",
				Clear:  false,
			},
			want: httpAction{
				Action: "http",
				Label:  "Nil URL",
				Method: "DELETE",
				Body:   "\"simple string body\"",
				Clear:  false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(&tt.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			want, err := json.Marshal(&tt.want)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}

			if !reflect.DeepEqual(got, want) {
				t.Errorf("json.Marshal() = %v, want %v", string(got), string(want))
			}
		})
	}
}

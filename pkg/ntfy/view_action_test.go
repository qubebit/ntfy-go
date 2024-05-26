package ntfy

import (
	"encoding/json"
	"net/url"
	"reflect"
	"testing"
)

func TestViewMarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		action  ViewAction
		want    viewAction
		wantErr bool
	}{
		{
			name: "Simple case",
			action: ViewAction{
				Label: "Test",
				Link:  &url.URL{Scheme: "https", Host: "example.com", Path: "/test"},
				Clear: true,
			},
			want: viewAction{
				Action: "view",
				Label:  "Test",
				URL:    "https://example.com/test",
				Clear:  true,
			},
			wantErr: false,
		},
		{
			name: "Empty URL",
			action: ViewAction{
				Label: "Empty URL",
				Clear: false,
			},
			want: viewAction{
				Action: "view",
				Label:  "Empty URL",
				Clear:  false,
			},
			wantErr: false,
		},
		{
			name: "Nil URL",
			action: ViewAction{
				Label: "Nil URL",
				Link:  nil,
				Clear: true,
			},
			want: viewAction{
				Action: "view",
				Label:  "Nil URL",
				Clear:  true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.action.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("ViewAction.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var gotStruct viewAction
			if err := json.Unmarshal(got, &gotStruct); err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			if !reflect.DeepEqual(gotStruct, tt.want) {
				t.Errorf("ViewAction.MarshalJSON() = %v, want %v", gotStruct, tt.want)
			}
		})
	}
}

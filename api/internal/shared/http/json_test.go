package sharedhttp

import (
	"strings"
	"testing"
)

func TestDecodeJSON(t *testing.T) {
	type request struct {
		Name string `json:"name"`
	}

	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{
			name: "accepts one JSON object with trailing whitespace",
			body: `{"name":"BankLab"} ` + "\n\t",
		},
		{
			name:    "rejects unknown field",
			body:    `{"name":"BankLab","extra":true}`,
			wantErr: true,
		},
		{
			name:    "rejects second JSON object",
			body:    `{"name":"BankLab"}{"name":"Other"}`,
			wantErr: true,
		},
		{
			name:    "rejects trailing JSON primitive",
			body:    `{"name":"BankLab"} true`,
			wantErr: true,
		},
		{
			name:    "rejects malformed JSON",
			body:    `{"name":`,
			wantErr: true,
		},
		{
			name:    "rejects empty body",
			body:    ``,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got request
			err := DecodeJSON(strings.NewReader(tt.body), &got)

			if tt.wantErr && err == nil {
				t.Fatal("expected decode error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no decode error, got %v", err)
			}
		})
	}
}

package json

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMarshalStruct(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		wantJSON string
	}{
		{
			name:     "SimpleMap",
			input:    map[string]string{"key": "value"},
			wantJSON: `"key": "value"`,
		},
		{
			name:     "Struct",
			input:    struct{ Name string }{Name: "test"},
			wantJSON: `"Name": "test"`,
		},
		{
			name:     "Nil",
			input:    nil,
			wantJSON: "null",
		},
		{
			name:     "EmptyStruct",
			input:    struct{}{},
			wantJSON: "{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MarshalStruct(tt.input)
			if !strings.Contains(string(got), tt.wantJSON) {
				t.Errorf("MarshalStruct() = %s, want to contain %s", got, tt.wantJSON)
			}
		})
	}
}

func TestWriteResult(t *testing.T) {
	tests := []struct {
		name           string
		payload        interface{}
		statusCode     int
		wantBody       string
		wantStatusCode int
		wantContentType string
	}{
		{
			name:            "OK",
			payload:         map[string]string{"status": "ok"},
			statusCode:      http.StatusOK,
			wantBody:        `"status":"ok"`,
			wantStatusCode:  http.StatusOK,
			wantContentType: "application/json",
		},
		{
			name:            "BadRequest",
			payload:         map[string]string{"error": "bad"},
			statusCode:      http.StatusBadRequest,
			wantBody:        `"error":"bad"`,
			wantStatusCode:  http.StatusBadRequest,
			wantContentType: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			WriteResult(tt.payload, rr, tt.statusCode)

			if rr.Code != tt.wantStatusCode {
				t.Errorf("WriteResult() status = %d, want %d", rr.Code, tt.wantStatusCode)
			}
			if ct := rr.Header().Get("Content-Type"); ct != tt.wantContentType {
				t.Errorf("WriteResult() Content-Type = %s, want %s", ct, tt.wantContentType)
			}
			if !strings.Contains(rr.Body.String(), tt.wantBody) {
				t.Errorf("WriteResult() body = %s, want to contain %s", rr.Body.String(), tt.wantBody)
			}
		})
	}
}

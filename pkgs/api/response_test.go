package api

import (
	"net/http"
	"testing"
)

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name         string
		errorMessage string
		statusCode   int
	}{
		{"BasicError", "something went wrong", http.StatusBadRequest},
		{"UnprocessableEntity", "invalid input", http.StatusUnprocessableEntity},
		{"InternalServerError", "server error", http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ErrorResponse(tt.errorMessage, tt.statusCode)
			if got.Success {
				t.Error("ErrorResponse() Success = true, want false")
			}
			if got.ErrorMessage != tt.errorMessage {
				t.Errorf("ErrorResponse() ErrorMessage = %v, want %v", got.ErrorMessage, tt.errorMessage)
			}
			if got.Data != nil {
				t.Errorf("ErrorResponse() Data = %v, want nil", got.Data)
			}
		})
	}
}

func TestSuccessResponse(t *testing.T) {
	tests := []struct {
		name    string
		payload interface{}
	}{
		{"StringPayload", "some data"},
		{"MapPayload", map[string]int{"count": 5}},
		{"NilPayload", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SuccessResponse(tt.payload)
			if !got.Success {
				t.Error("SuccessResponse() Success = false, want true")
			}
			if got.ErrorMessage != "" {
				t.Errorf("SuccessResponse() ErrorMessage = %v, want empty", got.ErrorMessage)
			}
		})
	}
}

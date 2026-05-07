package auth

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"Simple", "password123"},
		{"Special", "p@ssw0rd!#$"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HashPassword(tt.password)
			if h == tt.password {
				t.Error("HashPassword() returned plaintext password")
			}
			// bcrypt embeds a random salt, so two hashes of the same password differ —
			// verify via CheckPassword instead of string equality.
			if !CheckPassword(tt.password, h) {
				t.Error("HashPassword() produced a hash that does not verify")
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"Correct", "secret", true},
		{"Wrong", "wrong", false},
		{"Empty", "", false},
	}
	hash := HashPassword("secret")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPassword(tt.password, hash); got != tt.want {
				t.Errorf("CheckPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken(42, "REGULAR")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	if token == "" {
		t.Fatal("GenerateToken() returned empty token")
	}
	if !strings.Contains(token, ".") {
		t.Error("GenerateToken() token missing '.' separator")
	}
}

func TestValidateToken(t *testing.T) {
	token, err := GenerateToken(42, "REGULAR")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	payload, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}
	if payload.UserID != 42 {
		t.Errorf("ValidateToken() UserID = %v, want 42", payload.UserID)
	}
	if payload.UserType != "REGULAR" {
		t.Errorf("ValidateToken() UserType = %v, want REGULAR", payload.UserType)
	}
}

func TestValidateToken_Tampered(t *testing.T) {
	token, _ := GenerateToken(1, "REGULAR")
	parts := strings.SplitN(token, ".", 2)
	tampered := parts[0] + ".invalidsignature"

	_, err := ValidateToken(tampered)
	if err == nil {
		t.Error("ValidateToken() expected error for tampered token, got nil")
	}
}

func TestValidateToken_Malformed(t *testing.T) {
	_, err := ValidateToken("notavalidtoken")
	if err == nil {
		t.Error("ValidateToken() expected error for malformed token, got nil")
	}
}

package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"KnowWhoami/movie-ticketing/pkgs/auth"
	"KnowWhoami/movie-ticketing/pkgs/models"
)

func tokenFor(t *testing.T, userID int, userType string) string {
	t.Helper()
	tok, err := auth.GenerateToken(userID, userType)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	return tok
}

func okHandler(called *bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*called = true
		w.WriteHeader(http.StatusOK)
	})
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	called := false
	h := AuthMiddleware(okHandler(&called))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest("GET", "/", nil))

	if called {
		t.Error("next handler should not have been called")
	}
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	called := false
	h := AuthMiddleware(okHandler(&called))
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer notavalidtoken")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if called {
		t.Error("next handler should not have been called")
	}
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_ValidToken_CallsNext(t *testing.T) {
	called := false
	h := AuthMiddleware(okHandler(&called))
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenFor(t, 1, "REGULAR"))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if !called {
		t.Error("next handler should have been called")
	}
	if rr.Code != http.StatusOK {
		t.Errorf("want 200, got %d", rr.Code)
	}
}

func TestAuthMiddleware_StoresPayloadInContext(t *testing.T) {
	var gotPayload *auth.TokenPayload
	h := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPayload = TokenPayloadFromContext(r.Context())
	}))
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenFor(t, 42, "THEATRE_OWNER"))
	h.ServeHTTP(httptest.NewRecorder(), req)

	if gotPayload == nil {
		t.Fatal("TokenPayloadFromContext() returned nil")
	}
	if gotPayload.UserID != 42 {
		t.Errorf("UserID = %d, want 42", gotPayload.UserID)
	}
	if gotPayload.UserType != "THEATRE_OWNER" {
		t.Errorf("UserType = %s, want THEATRE_OWNER", gotPayload.UserType)
	}
}

func TestRequireRole_AllowedRole(t *testing.T) {
	called := false
	h := RequireRole(models.UserTypeTheatreOwner)(okHandler(&called))
	req := httptest.NewRequest("POST", "/cinema", nil)
	req.Header.Set("Authorization", "Bearer "+tokenFor(t, 1, "THEATRE_OWNER"))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if !called {
		t.Error("next handler should have been called for THEATRE_OWNER")
	}
	if rr.Code != http.StatusOK {
		t.Errorf("want 200, got %d", rr.Code)
	}
}

func TestRequireRole_WrongRole(t *testing.T) {
	called := false
	h := RequireRole(models.UserTypeTheatreOwner)(okHandler(&called))
	req := httptest.NewRequest("POST", "/cinema", nil)
	req.Header.Set("Authorization", "Bearer "+tokenFor(t, 1, "REGULAR"))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if called {
		t.Error("next handler should not have been called for REGULAR user on admin route")
	}
	if rr.Code != http.StatusForbidden {
		t.Errorf("want 403, got %d", rr.Code)
	}
}

func TestRequireRole_NoToken(t *testing.T) {
	called := false
	h := RequireRole(models.UserTypeTheatreOwner)(okHandler(&called))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest("POST", "/cinema", nil))

	if called {
		t.Error("next handler should not have been called without token")
	}
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", rr.Code)
	}
}

func TestTokenPayloadFromContext_Empty(t *testing.T) {
	if got := TokenPayloadFromContext(context.Background()); got != nil {
		t.Errorf("want nil for empty context, got %v", got)
	}
}

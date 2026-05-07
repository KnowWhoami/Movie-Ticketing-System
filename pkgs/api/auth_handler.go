package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	jsonHelper "KnowWhoami/movie-ticketing/internal/json"
	"KnowWhoami/movie-ticketing/pkgs/auth"
	"KnowWhoami/movie-ticketing/pkgs/models"
)

type contextKey int

const tokenPayloadKey contextKey = 0

type registerInput struct {
	Name     string          `json:"name"`
	Email    string          `json:"email"`
	Password string          `json:"password"`
	UserType models.UserType `json:"user_type"`
}

type loginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authOutput struct {
	Token string `json:"token"`
	User  struct {
		ID       int             `json:"id"`
		Name     string          `json:"name"`
		Email    string          `json:"email"`
		UserType models.UserType `json:"user_type"`
	} `json:"user"`
}

func (h *Handler) Register(request *http.Request, writer http.ResponseWriter) {
	var inp registerInput
	if err := json.NewDecoder(request.Body).Decode(&inp); err != nil {
		resp := ErrorResponse("invalid request body", http.StatusBadRequest)
		jsonHelper.WriteResult(&resp, writer, http.StatusBadRequest)
		return
	}

	if inp.Name == "" || inp.Email == "" || inp.Password == "" {
		resp := ErrorResponse("name, email, and password are required", http.StatusBadRequest)
		jsonHelper.WriteResult(&resp, writer, http.StatusBadRequest)
		return
	}

	if inp.UserType == "" {
		inp.UserType = models.UserTypeRegular
	}

	user := models.User{
		Name:     inp.Name,
		Email:    inp.Email,
		Password: auth.HashPassword(inp.Password),
		UserType: inp.UserType,
	}

	if result := h.db.Create(&user); result.Error != nil {
		_ = h.logger.Log("handler", "Register", "err", result.Error)
		resp := ErrorResponse(result.Error.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}

	token, err := auth.GenerateToken(user.ID, string(user.UserType))
	if err != nil {
		_ = h.logger.Log("handler", "Register", "err", err)
		resp := ErrorResponse("failed to generate token", http.StatusInternalServerError)
		jsonHelper.WriteResult(&resp, writer, http.StatusInternalServerError)
		return
	}

	_ = h.logger.Log("handler", "Register", "user_id", user.ID, "user_type", user.UserType)
	var out authOutput
	out.Token = token
	out.User.ID = user.ID
	out.User.Name = user.Name
	out.User.Email = user.Email
	out.User.UserType = user.UserType

	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusCreated)
}

func (h *Handler) Login(request *http.Request, writer http.ResponseWriter) {
	var inp loginInput
	if err := json.NewDecoder(request.Body).Decode(&inp); err != nil {
		resp := ErrorResponse("invalid request body", http.StatusBadRequest)
		jsonHelper.WriteResult(&resp, writer, http.StatusBadRequest)
		return
	}

	if inp.Email == "" || inp.Password == "" {
		resp := ErrorResponse("email and password are required", http.StatusBadRequest)
		jsonHelper.WriteResult(&resp, writer, http.StatusBadRequest)
		return
	}

	var user models.User
	if result := h.db.Where("email = ?", inp.Email).First(&user); result.Error != nil {
		_ = h.logger.Log("handler", "Login", "err", result.Error)
		resp := ErrorResponse("invalid email or password", http.StatusUnauthorized)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnauthorized)
		return
	}

	if !auth.CheckPassword(inp.Password, user.Password) {
		_ = h.logger.Log("handler", "Login", "err", "password mismatch", "user_id", user.ID)
		resp := ErrorResponse("invalid email or password", http.StatusUnauthorized)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID, string(user.UserType))
	if err != nil {
		_ = h.logger.Log("handler", "Login", "err", err)
		resp := ErrorResponse("failed to generate token", http.StatusInternalServerError)
		jsonHelper.WriteResult(&resp, writer, http.StatusInternalServerError)
		return
	}

	_ = h.logger.Log("handler", "Login", "user_id", user.ID, "user_type", user.UserType)
	var out authOutput
	out.Token = token
	out.User.ID = user.ID
	out.User.Name = user.Name
	out.User.Email = user.Email
	out.User.UserType = user.UserType

	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

// AuthMiddleware validates the Bearer token and stores the payload in the request context.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		payload, ok := extractToken(writer, request)
		if !ok {
			return
		}
		next.ServeHTTP(writer, request.WithContext(
			context.WithValue(request.Context(), tokenPayloadKey, payload),
		))
	})
}

// RequireRole returns a middleware that only allows users with one of the given roles.
// It builds on AuthMiddleware — the token is validated first, then the role is checked.
func RequireRole(roles ...models.UserType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return AuthMiddleware(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			payload := TokenPayloadFromContext(request.Context())
			for _, role := range roles {
				if models.UserType(payload.UserType) == role {
					next.ServeHTTP(writer, request)
					return
				}
			}
			resp := ErrorResponse("insufficient permissions", http.StatusForbidden)
			jsonHelper.WriteResult(&resp, writer, http.StatusForbidden)
		}))
	}
}

// TokenPayloadFromContext extracts the token payload stored by AuthMiddleware.
func TokenPayloadFromContext(ctx context.Context) *auth.TokenPayload {
	p, _ := ctx.Value(tokenPayloadKey).(*auth.TokenPayload)
	return p
}

// extractToken parses the Bearer token from the Authorization header.
// Returns the payload and true on success; writes an error response and returns false on failure.
func extractToken(writer http.ResponseWriter, request *http.Request) (*auth.TokenPayload, bool) {
	header := request.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		resp := ErrorResponse("missing or invalid authorization header", http.StatusUnauthorized)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnauthorized)
		return nil, false
	}
	payload, err := auth.ValidateToken(strings.TrimPrefix(header, "Bearer "))
	if err != nil {
		resp := ErrorResponse("invalid or expired token", http.StatusUnauthorized)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnauthorized)
		return nil, false
	}
	return payload, true
}

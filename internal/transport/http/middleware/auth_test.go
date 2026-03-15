package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	secret := "test-secret"
	userID := uint(123)

	generateToken := func(id uint, secret string, expired bool) string {
		claims := jwt.MapClaims{
			"sub":  float64(id),
			"role": "user",
			"exp":  time.Now().Add(time.Hour).Unix(),
		}
		if expired {
			claims["exp"] = time.Now().Add(-time.Hour).Unix()
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))
		return tokenString
	}

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "Valid Token",
			authHeader:     "Bearer " + generateToken(userID, secret, false),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing Authorization Header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Format",
			authHeader:     "InvalidFormat",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Wrong Secret",
			authHeader:     "Bearer " + generateToken(userID, "wrong-secret", false),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Expired Token",
			authHeader:     "Bearer " + generateToken(userID, secret, true),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := Auth(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				val := r.Context().Value(UserIDKey)
				assert.Equal(t, userID, val)
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name          string
		header        string
		expectedToken string
		expectedErr   bool
	}{
		{"Valid Bearer", "Bearer mytoken", "mytoken", false},
		{"Empty Header", "", "", true},
		{"Wrong Format", "Token mytoken", "", true},
		{"No Token", "Bearer ", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			token, err := extractToken(req)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, token)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	secret := "secret"

	t.Run("Valid Token", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": float64(1), "role": "user"})
		ts, _ := token.SignedString([]byte(secret))

		id, _, err := validateToken(ts, secret)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), id)
	})

	t.Run("Invalid Signing Method", func(t *testing.T) {
		// This is a simplified check, jwt.Parse handles this
		ts := fmt.Sprintf("%s.%s.", "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0", "eyJzdWIiOjF9")

		_, _, err := validateToken(ts, secret)
		assert.Error(t, err)
	})
}

package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/m-milek/leszmonitor/auth"
	"github.com/m-milek/leszmonitor/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create a valid JWT token
func createTestToken(t *testing.T, claims *auth.UserClaims, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	require.NoError(t, err)
	return tokenString
}

// Helper function to create a test handler that captures the context
func createTestHandler(t *testing.T) (http.Handler, **auth.UserClaims) {
	var capturedClaims *auth.UserClaims

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetUserFromContext(r.Context())
		if ok {
			capturedClaims = claims
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	return handler, &capturedClaims
}

func TestJwtAuth(t *testing.T) {
	// Set up test JWT secret
	testSecret := "test-secret-key"
	os.Setenv(config.JwtSecret, testSecret)
	defer os.Unsetenv(config.JwtSecret)

	tests := []struct {
		name           string
		setupAuth      func() string
		setupEnv       func()
		expectedStatus int
		expectedError  string
		expectClaims   bool
	}{
		{
			name: "valid token with Bearer prefix",
			setupAuth: func() string {
				claims := &auth.UserClaims{
					Username: "testuser",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					},
				}
				token := createTestToken(t, claims, testSecret)
				return "Bearer " + token
			},
			expectedStatus: http.StatusOK,
			expectClaims:   true,
		},
		{
			name: "missing authorization header",
			setupAuth: func() string {
				return ""
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Unauthorized: No token provided",
		},
		{
			name:           "invalid token returns generic unauthorized",
			setupAuth:      func() string { return "Bearer invalid.token.value" },
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Unauthorized: Invalid token",
		},
		{
			name: "missing JWT secret returns unauthorized",
			setupAuth: func() string {
				return "Bearer some.token.here"
			},
			setupEnv: func() {
				os.Unsetenv(config.JwtSecret)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Unauthorized: Invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset environment
			os.Setenv(config.JwtSecret, testSecret)

			if tt.setupEnv != nil {
				tt.setupEnv()
			}

			// Create test request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if auth := tt.setupAuth(); auth != "" {
				req.Header.Set("Authorization", auth)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Create test handler
			testHandler, capturedClaims := createTestHandler(t)

			// Apply middleware
			handler := JwtAuth(testHandler)

			// Execute request
			handler.ServeHTTP(rr, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Assert error message if expected
			if tt.expectedError != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedError)
			}

			// Assert claims were set in context
			if tt.expectClaims {
				require.NotNil(t, *capturedClaims)
				assert.Equal(t, "testuser", (*capturedClaims).Username)
			} else {
				assert.Nil(t, *capturedClaims)
			}
		})
	}
}

func TestSetUserContext(t *testing.T) {
	ctx := context.Background()
	claims := &auth.UserClaims{
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "user123",
		},
	}

	newCtx := SetUserContext(ctx, claims)

	// Verify context contains the claims
	value := newCtx.Value(userClaimsKey)
	assert.NotNil(t, value)

	retrievedClaims, ok := value.(*auth.UserClaims)
	assert.True(t, ok)
	assert.Equal(t, claims.Username, retrievedClaims.Username)
	assert.Equal(t, claims.Subject, retrievedClaims.Subject)
}

func TestGetUserFromContext(t *testing.T) {
	t.Run("context with claims", func(t *testing.T) {
		ctx := context.Background()
		expectedClaims := &auth.UserClaims{
			Username: "testuser",
		}

		ctx = SetUserContext(ctx, expectedClaims)

		claims, ok := GetUserFromContext(ctx)
		assert.True(t, ok)
		assert.NotNil(t, claims)
		assert.Equal(t, expectedClaims.Username, claims.Username)
	})

	t.Run("context without claims", func(t *testing.T) {
		ctx := context.Background()

		claims, ok := GetUserFromContext(ctx)
		assert.False(t, ok)
		assert.Nil(t, claims)
	})

	t.Run("context with wrong type", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, userClaimsKey, "not-a-claim")

		claims, ok := GetUserFromContext(ctx)
		assert.False(t, ok)
		assert.Nil(t, claims)
	})
}

func TestProjectAuthFromRequest(t *testing.T) {
	tests := []struct {
		name          string
		setupRequest  func() *http.Request
		expectedAuth  *ProjectAuth
		expectedError string
	}{
		{
			name: "valid request with project slug and user claims",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/projects/proj123", nil)
				req.SetPathValue("projectId", "proj123")

				claims := &auth.UserClaims{Username: "testuser"}
				ctx := SetUserContext(req.Context(), claims)
				return req.WithContext(ctx)
			},
			expectedAuth: &ProjectAuth{
				ProjectID: "proj123",
				Username:  "testuser",
			},
		},
		{
			name: "missing project slug",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/projects/", nil)
				// No projectId path value set

				claims := &auth.UserClaims{Username: "testuser"}
				ctx := SetUserContext(req.Context(), claims)
				return req.WithContext(ctx)
			},
			expectedError: "projectID is required",
		},
		{
			name: "missing user claims in context",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/projects/proj123", nil)
				req.SetPathValue("projectId", "proj123")
				// No user claims in context
				return req
			},
			expectedError: "user claims not found in context",
		},
		{
			name: "empty username in claims",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/projects/proj123", nil)
				req.SetPathValue("projectId", "proj123")

				claims := &auth.UserClaims{Username: ""} // Empty username
				ctx := SetUserContext(req.Context(), claims)
				return req.WithContext(ctx)
			},
			expectedError: "username is missing in user claims",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setupRequest()

			auth, err := ProjectAuthFromRequest(req)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, auth)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, auth)
				assert.Equal(t, tt.expectedAuth.ProjectID, auth.ProjectID)
				assert.Equal(t, tt.expectedAuth.Username, auth.Username)
			}
		})
	}
}

// Test the responseWriter wrapper if it's used
func TestResponseWriterWrapper(t *testing.T) {
	// This assumes you have a newResponseWriter function that wraps http.ResponseWriter
	// If not, you can skip this test

	t.Run("wrapper passes through writes", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		wrapper := newResponseWriter(recorder)

		// Test writing header
		wrapper.WriteHeader(http.StatusCreated)
		assert.Equal(t, http.StatusCreated, recorder.Code)

		// Test writing body
		n, err := wrapper.Write([]byte("test response"))
		assert.NoError(t, err)
		assert.Equal(t, 13, n)
		assert.Equal(t, "test response", recorder.Body.String())
	})
}

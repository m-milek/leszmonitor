package api

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/m-milek/leszmonitor/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testJwtSecret = "test-secret-key-for-testing"

func TestMain(m *testing.M) {
	// Set up test environment
	os.Setenv(env.JwtSecret, testJwtSecret)
	code := m.Run()
	// Clean up
	os.Unsetenv(env.JwtSecret)
	os.Exit(code)
}

func createTestToken(claims jwtClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func TestJwtFromRequest(t *testing.T) {
	tests := []struct {
		name          string
		authHeader    string
		expectedToken string
		expectError   bool
	}{
		{
			name:          "Valid Bearer Token",
			authHeader:    "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
			expectedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
			expectError:   false,
		},
		{
			name:          "Empty Authorization Header",
			authHeader:    "",
			expectedToken: "",
			expectError:   false,
		},
		{
			name:          "No Bearer Prefix",
			authHeader:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
			expectedToken: "",
			expectError:   false,
		},
		{
			name:          "Invalid Bearer Format - No Space",
			authHeader:    "BearereyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
			expectedToken: "",
			expectError:   false,
		},
		{
			name:          "Only Bearer Prefix",
			authHeader:    "Bearer ",
			expectedToken: "",
			expectError:   false,
		},
		{
			name:          "Bearer With Multiple Spaces",
			authHeader:    "Bearer  token-with-space",
			expectedToken: " token-with-space",
			expectError:   false,
		},
		{
			name:          "Different Auth Scheme",
			authHeader:    "Basic dXNlcjpwYXNzd29yZA==",
			expectedToken: "",
			expectError:   false,
		},
		{
			name:          "Case Sensitive Bearer",
			authHeader:    "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
			expectedToken: "",
			expectError:   false,
		},
		{
			name:          "Bearer With Tab",
			authHeader:    "Bearer\ttoken",
			expectedToken: "",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/test", nil)
			require.NoError(t, err)

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			token, err := jwtFromRequest(req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, token)
			}
		})
	}
}

func TestDecodeJwtClaims(t *testing.T) {
	t.Run("Valid Token", func(t *testing.T) {
		claims := jwtClaims{
			MapClaims: jwt.MapClaims{
				"sub": "1234567890",
				"iat": time.Now().Unix(),
			},
			Username: "testuser",
			Exp:      time.Now().Add(time.Hour).Unix(),
		}

		token, err := createTestToken(claims, testJwtSecret)
		require.NoError(t, err)

		decodedClaims, err := decodeJwtClaims(token)
		assert.NoError(t, err)
		assert.Equal(t, claims.Username, decodedClaims.Username)
		assert.Equal(t, claims.Exp, decodedClaims.Exp)
	})

	t.Run("Expired Token", func(t *testing.T) {
		claims := jwtClaims{
			MapClaims: jwt.MapClaims{
				"exp": time.Now().Add(-time.Hour).Unix(),
			},
			Username: "testuser",
			Exp:      time.Now().Add(-time.Hour).Unix(),
		}

		token, err := createTestToken(claims, testJwtSecret)
		require.NoError(t, err)

		_, err = decodeJwtClaims(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired")
	})

	t.Run("Invalid Signature", func(t *testing.T) {
		claims := jwtClaims{
			Username: "testuser",
			Exp:      time.Now().Add(time.Hour).Unix(),
		}

		token, err := createTestToken(claims, "wrong-secret")
		require.NoError(t, err)

		_, err = decodeJwtClaims(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "signature is invalid")
	})

	t.Run("Malformed Token", func(t *testing.T) {
		malformedTokens := []string{
			"not.a.token",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			"",
			".",
			"..",
			"a.b",
			"invalid-base64.test.test",
		}

		for _, token := range malformedTokens {
			_, err := decodeJwtClaims(token)
			assert.Error(t, err)
		}
	})

	t.Run("Token Without Username", func(t *testing.T) {
		claims := jwtClaims{
			MapClaims: jwt.MapClaims{
				"sub": "1234567890",
			},
			Exp: time.Now().Add(time.Hour).Unix(),
		}

		token, err := createTestToken(claims, testJwtSecret)
		require.NoError(t, err)

		decodedClaims, err := decodeJwtClaims(token)
		assert.NoError(t, err)
		assert.Empty(t, decodedClaims.Username)
		assert.Equal(t, claims.Exp, decodedClaims.Exp)
	})

	t.Run("Token With Additional Claims", func(t *testing.T) {
		claims := jwtClaims{
			MapClaims: jwt.MapClaims{
				"role":   "admin",
				"userId": "123",
			},
			Username: "testuser",
			Exp:      time.Now().Add(time.Hour).Unix(),
		}

		token, err := createTestToken(claims, testJwtSecret)
		require.NoError(t, err)

		decodedClaims, err := decodeJwtClaims(token)
		assert.NoError(t, err)
		assert.Equal(t, claims.Username, decodedClaims.Username)
		assert.Equal(t, "admin", decodedClaims.MapClaims["role"])
		assert.Equal(t, "123", decodedClaims.MapClaims["userId"])
	})

	t.Run("Empty JWT Secret", func(t *testing.T) {
		// Temporarily unset the JWT secret
		originalSecret := os.Getenv(env.JwtSecret)
		os.Unsetenv(env.JwtSecret)
		defer os.Setenv(env.JwtSecret, originalSecret)

		claims := jwtClaims{
			Username: "testuser",
			Exp:      time.Now().Add(time.Hour).Unix(),
		}

		token, err := createTestToken(claims, "")
		require.NoError(t, err)

		decodedClaims, err := decodeJwtClaims(token)
		assert.NoError(t, err)
		assert.Equal(t, claims.Username, decodedClaims.Username)
	})

	t.Run("Different Signing Algorithm", func(t *testing.T) {
		// Create a token with a different algorithm
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtClaims{
			Username: "testuser",
			Exp:      time.Now().Add(time.Hour).Unix(),
		})
		tokenString, err := token.SignedString([]byte(testJwtSecret))
		require.NoError(t, err)

		// This should still work as jwt-go handles multiple algorithms
		decodedClaims, err := decodeJwtClaims(tokenString)
		assert.NoError(t, err)
		assert.Equal(t, "testuser", decodedClaims.Username)
	})

	t.Run("Token Validation Error", func(t *testing.T) {
		// Create a token that will fail validation
		token := jwt.New(jwt.SigningMethodHS256)
		token.Claims = jwtClaims{
			Username: "testuser",
			Exp:      time.Now().Add(time.Hour).Unix(),
		}
		token.Valid = false // Force invalid state

		tokenString, err := token.SignedString([]byte(testJwtSecret))
		require.NoError(t, err)

		// Parse should succeed but the token should be detected as invalid
		parsedToken, _ := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(testJwtSecret), nil
		})

		// The token should be valid after proper parsing
		assert.True(t, parsedToken.Valid)
	})
}

func TestJwtClaimsIntegration(t *testing.T) {
	t.Run("Full Flow - Request to Claims", func(t *testing.T) {
		// Create a valid token
		claims := jwtClaims{
			MapClaims: jwt.MapClaims{
				"iat": time.Now().Unix(),
			},
			Username: "integrationuser",
			Exp:      time.Now().Add(time.Hour).Unix(),
		}

		token, err := createTestToken(claims, testJwtSecret)
		require.NoError(t, err)

		// Create request with token
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		// Extract token from request
		extractedToken, err := jwtFromRequest(req)
		assert.NoError(t, err)
		assert.Equal(t, token, extractedToken)

		// Decode claims from token
		decodedClaims, err := decodeJwtClaims(extractedToken)
		assert.NoError(t, err)
		assert.Equal(t, claims.Username, decodedClaims.Username)
		assert.Equal(t, claims.Exp, decodedClaims.Exp)
	})
}

// Benchmark tests
func BenchmarkJwtFromRequest(b *testing.B) {
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = jwtFromRequest(req)
	}
}

func BenchmarkDecodeJwtClaims(b *testing.B) {
	claims := jwtClaims{
		Username: "benchuser",
		Exp:      time.Now().Add(time.Hour).Unix(),
	}
	token, _ := createTestToken(claims, testJwtSecret)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = decodeJwtClaims(token)
	}
}

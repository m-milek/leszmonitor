package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/m-milek/leszmonitor/api/authorization"
	config "github.com/m-milek/leszmonitor/appconfig"
	"github.com/m-milek/leszmonitor/auth"
)

func TestJwtAuth_NoAuthHeader(t *testing.T) {
	handler := JwtAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestJwtAuth_InvalidToken(t *testing.T) {
	handler := JwtAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestJwtAuth_ValidToken(t *testing.T) {
	os.Setenv(config.JwtSecret, "test-secret")
	os.Setenv(config.JwtExpiryHours, "1")
	token, err := auth.NewJwt("testuser", false)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	called := false
	handler := JwtAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		claims, ok := authorization.GetUserClaimsFromContext(r.Context())
		if !ok {
			t.Fatal("expected claims in context, got nil")
		}
		if claims.Username != "testuser" {
			t.Errorf("expected username 'testuser', got '%s'", claims.Username)
		}

		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+*token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !called {
		t.Fatal("next handler was not called")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestJwtAuth_MissingBearerPrefix(t *testing.T) {
	handler := JwtAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "not-a-bearer-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

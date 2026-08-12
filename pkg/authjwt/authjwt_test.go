package authjwt

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestNewTokenAndParse(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-that-is-long-enough-123456")
	raw, err := NewToken("user-1", "user@example.com", "customer", time.Now())
	if err != nil {
		t.Fatalf("NewToken() error = %v", err)
	}
	claims, err := Parse(raw)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if claims.UserID != "user-1" || claims.Role != "customer" {
		t.Fatalf("unexpected claims: %#v", claims)
	}
}

func TestParseRejectsMissingSecretAndWrongAlgorithm(t *testing.T) {
	t.Setenv("JWT_SECRET", "")
	if _, err := Parse("anything"); err != ErrSecretNotConfigured {
		t.Fatalf("Parse() missing secret error = %v", err)
	}

	if err := os.Setenv("JWT_SECRET", "test-secret-that-is-long-enough-123456"); err != nil {
		t.Fatal(err)
	}
	claims := Claims{UserID: "user-1", Email: "user@example.com", Role: "customer", RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))}}
	raw, err := jwt.NewWithClaims(jwt.SigningMethodHS384, claims).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Parse(raw); err == nil {
		t.Fatal("Parse() accepted HS384 token")
	}
}

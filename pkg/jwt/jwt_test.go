package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestInit(t *testing.T) {
	Init("test-secret-key-12345")
	if len(jwtSecret) == 0 {
		t.Error("Init() failed to set jwtSecret")
	}
}

func TestGenerateAccessToken(t *testing.T) {
	Init("test-secret-key-12345")

	token, err := GenerateAccessToken("user123")
	if err != nil {
		t.Errorf("GenerateAccessToken() error = %v", err)
	}
	if token == "" {
		t.Error("GenerateAccessToken() returned empty token")
	}
}

func TestParseToken(t *testing.T) {
	Init("test-secret-key-12345")

	token, err := GenerateAccessToken("user123")
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Errorf("ParseToken() error = %v", err)
	}
	if claims.UserID != "user123" {
		t.Errorf("ParseToken() UserID = %v, want %v", claims.UserID, "user123")
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	Init("test-secret-key-12345")

	_, err := ParseToken("invalid-token-string")
	if err == nil {
		t.Error("ParseToken() should return error for invalid token")
	}
}

func TestParseToken_ExpiredToken(t *testing.T) {
	Init("test-secret-key-12345")

	now := time.Now()
	claims := Claims{
		UserID: "user123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	_, err = ParseToken(tokenString)
	if err == nil {
		t.Error("ParseToken() should return error for expired token")
	}
}

func TestParseTokenAndGetUserID(t *testing.T) {
	Init("test-secret-key-12345")

	token, err := GenerateAccessToken("user456")
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	userID, err := ParseTokenAndGetUserID(token)
	if err != nil {
		t.Errorf("ParseTokenAndGetUserID() error = %v", err)
	}
	if userID != "user456" {
		t.Errorf("ParseTokenAndGetUserID() = %v, want %v", userID, "user456")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	token, err := GenerateRefreshToken()
	if err != nil {
		t.Errorf("GenerateRefreshToken() error = %v", err)
	}
	if len(token) != 64 {
		t.Errorf("GenerateRefreshToken() length = %v, want 64", len(token))
	}
}

func TestGenerateRefreshToken_Unique(t *testing.T) {
	token1, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	token2, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	if token1 == token2 {
		t.Error("GenerateRefreshToken() should return unique tokens")
	}
}

package jwtutil

import (
	"testing"
	"time"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret := "unit-test-secret"
	token, err := GenerateToken(secret, "alice", "USER", time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	claims, err := ParseToken(secret, token)
	if err != nil {
		t.Fatalf("ParseToken returned error: %v", err)
	}

	if claims.Username != "alice" {
		t.Fatalf("expected username alice, got %s", claims.Username)
	}
	if claims.Role != "USER" {
		t.Fatalf("expected role USER, got %s", claims.Role)
	}
}


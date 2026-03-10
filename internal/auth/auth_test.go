package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "secretpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("Error checking password: %v", err)
	}
	if !match {
		t.Fatal("Password should match hash")
	}
}

func TestCheckPasswordHashWrong(t *testing.T) {
	hash, err := HashPassword("correctpassword")
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	match, err := CheckPasswordHash("wrongpassword", hash)
	if err != nil {
		t.Fatalf("Error checking password: %v", err)
	}
	if match {
		t.Fatal("Wrong password should not match hash")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("Error making JWT: %v", err)
	}

	returnedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Error validating JWT: %v", err)
	}

	if returnedID != userID {
		t.Fatalf("Expected user ID %v, got %v", userID, returnedID)
	}
}

func TestExpiredJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"

	token, err := MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("Error making JWT: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("Expired token should fail validation")
	}
}

func TestWrongSecretJWT(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, "correct-secret", time.Hour)
	if err != nil {
		t.Fatalf("Error making JWT: %v", err)
	}

	_, err = ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Fatal("Token signed with wrong secret should fail validation")
	}
}

func TestHashPasswordDifferentHashes(t *testing.T) {
	password := "samepassword"
	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	if hash1 == hash2 {
		t.Fatal("Same password should produce different hashes due to unique salts")
	}
}

func TestValidateJWTEmptyToken(t *testing.T) {
	_, err := ValidateJWT("", "secret")
	if err == nil {
		t.Fatal("Empty token should fail validation")
	}
}

func TestValidateJWTGarbageToken(t *testing.T) {
	_, err := ValidateJWT("not.a.real.token", "secret")
	if err == nil {
		t.Fatal("Garbage token should fail validation")
	}
}

func TestDifferentUsersGetDifferentTokens(t *testing.T) {
	secret := "test-secret"
	user1 := uuid.New()
	user2 := uuid.New()

	token1, err := MakeJWT(user1, secret, time.Hour)
	if err != nil {
		t.Fatalf("Error making JWT: %v", err)
	}

	token2, err := MakeJWT(user2, secret, time.Hour)
	if err != nil {
		t.Fatalf("Error making JWT: %v", err)
	}

	if token1 == token2 {
		t.Fatal("Different users should get different tokens")
	}

	returnedID1, err := ValidateJWT(token1, secret)
	if err != nil {
		t.Fatalf("Error validating token1: %v", err)
	}

	returnedID2, err := ValidateJWT(token2, secret)
	if err != nil {
		t.Fatalf("Error validating token2: %v", err)
	}

	if returnedID1 != user1 {
		t.Fatalf("Token1 should return user1 ID, got %v", returnedID1)
	}

	if returnedID2 != user2 {
		t.Fatalf("Token2 should return user2 ID, got %v", returnedID2)
	}
}
func TestGetBearerTokenValid(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer my-token-123")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("Error getting bearer token: %v", err)
	}

	if token != "my-token-123" {
		t.Fatalf("Expected 'my-token-123', got '%s'", token)
	}
}

func TestGetBearerTokenMissing(t *testing.T) {
	headers := http.Header{}

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("Should error when Authorization header is missing")
	}
}

func TestGetBearerTokenNoBearer(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Basic my-token-123")

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("Should error when Authorization header doesn't start with Bearer")
	}
}

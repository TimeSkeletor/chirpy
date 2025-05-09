package auth

import (
	"testing"
)

// Test Password Hashing
func TestHashing(t *testing.T) {
	password := "MyPassword"
	hashed, _ := HashPassword(password)

	err := CheckPassword(string(hashed), password)

	if err != nil {
		t.Errorf("Password check failed: %v", err)
	}
}
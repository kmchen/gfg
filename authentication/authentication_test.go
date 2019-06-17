package authentication

import (
	"testing"
)

func TestValidToken(t *testing.T) {
	auth := &Authentication{}
	auth.Populate()
	clientId := "clientId0"
	token, err := auth.GenToken(clientId)
	if err != nil {
		t.Errorf("GenToken got an error, got: %v", err)
	}
	isValidToken, err := auth.IsValidToken(token)
	if err != nil {
		t.Errorf("ParseToken got an error, got: %v", err)
	}
	if !isValidToken {
		t.Errorf("Token %s is not valid", token)
	}
}

func TestInValidToken(t *testing.T) {
	auth := &Authentication{}
	auth.Populate()
	clientId := "clientId1"
	token, err := auth.GenToken(clientId)
	if err != nil {
		t.Errorf("GenToken got an error, got: %v", err)
	}
	isValidToken, err := auth.IsValidToken(token)
	if err != nil {
		t.Errorf("ParseToken got an error, got: %v", err)
	}
	if isValidToken {
		t.Errorf("Token %s should be invalid", token)
	}
}

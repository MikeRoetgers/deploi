package main

import (
	"testing"

	"github.com/MikeRoetgers/deploi/protobuf"
	jwt "github.com/dgrijalva/jwt-go"
)

func TestAuthEmptyHeader(t *testing.T) {
	_, pbErr := checkAuthentication(nil)
	if pbErr == nil {
		t.Error("Expected auth error but auth check found no problem")
	}
}

func TestAuthEmptyToken(t *testing.T) {
	_, pbErr := checkAuthentication(&protobuf.RequestHeader{
		Token: "",
	})
	if pbErr == nil {
		t.Error("Expected auth error but auth check found no problem")
	}
}

func TestAuthInvalidToken(t *testing.T) {
	_, pbErr := checkAuthentication(&protobuf.RequestHeader{
		Token: "CLEARLY_INVALID_TOKEN",
	})
	if pbErr == nil {
		t.Error("Expected auth error but auth check found no problem")
	}
}

func TestAuthValidToken(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &DeploidClaims{})
	signed, err := token.SignedString(JWTKey)
	if err != nil {
		t.Fatalf("Failed to sign test JWT token: %s", err)
	}
	_, pbErr := checkAuthentication(&protobuf.RequestHeader{
		Token: signed,
	})
	if pbErr != nil {
		t.Errorf("Expected positive auth result but got: %s", err)
	}
}

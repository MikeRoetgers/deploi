package main

import (
	"github.com/MikeRoetgers/deploi/protobuf"
	jwt "github.com/dgrijalva/jwt-go"
)

type DeploidClaims struct {
	User string
	jwt.StandardClaims
}

func checkAuthentication(header *protobuf.RequestHeader) (*DeploidClaims, *protobuf.Error) {
	if header == nil {
		return nil, &protobuf.Error{
			Code:    "REQUEST_INVALID",
			Message: "The received request does not include the mandatory RequestHeader.",
		}
	}
	if header.Token == "" {
		return nil, &protobuf.Error{
			Code:    "TOKEN_MISSING",
			Message: "Please log in first before you try to access protected endpoints.",
		}
	}
	token, err := jwt.ParseWithClaims(header.Token, &DeploidClaims{}, func(*jwt.Token) (interface{}, error) {
		return JWTKey, nil
	})
	if err != nil {
		log.Errorf("Failed to parse token: %s", err)
		return nil, &protobuf.Error{
			Code:    "TOKEN_INVALID",
			Message: "The provided token is invalid. Please log in again.",
		}
	}
	var claims *DeploidClaims
	var ok bool
	if claims, ok = token.Claims.(*DeploidClaims); !ok || !token.Valid {
		return nil, &protobuf.Error{
			Code:    "TOKEN_INVALID",
			Message: "The provided token is invalid. Please log in again. (Claims)",
		}
	}
	return claims, nil
}

package token

import (
	"errors"
	"github.com/golang-jwt/jwt"
	errors2 "ricardo/party-service/pkg/errors"
)

//import (
//	"errors"
//	"github.com/golang-jwt/jwt"
//)

const (
	AuthorizationHeader = "Authorization"
	BearerType          = "Bearer "
)

func ExtractTokenFromHeader(token string) (string, error) {
	if len(token) <= len(BearerType) {
		return "", errors.New("access token format is invalid")
	}
	return token[len(BearerType):], nil
}

func Parse(tokenStr string, key []byte) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New(errors2.InvalidToken)
		}
		return key, nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return token, nil
}

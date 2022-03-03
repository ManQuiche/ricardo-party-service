package token

//import (
//	"errors"
//	"github.com/golang-jwt/jwt"
//)

const (
	AuthorizationHeader = "Authorization"
	BearerType          = "Bearer "
)

//func ExtractTokenFromHeader(token string) (string, error) {
//	if len(token) <= len(BearerType) {
//		return "", errors.New("access token format is invalid")
//	}
//	return token[len(BearerType):], nil
//}
//
//func ExtractSubFromToken(tokenStr string) uint {
//	token, err := jwt.Parse
//
//
//	if claims, ok := token.
//}

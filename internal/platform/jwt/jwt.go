package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type accessTokenClaim struct {
	ID      uint
	Email   string
	KeyType string
	jwt.StandardClaims
}

type AccessTokenInfo struct {
	ID    uint
	Email string
}

type JWT struct {
	secretKey string
}

func New() *JWT {
	return &JWT{
		secretKey: os.Getenv("SECRET"),
	}
}

func (s *JWT) GenerateToken(userID uint, email string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute).Unix()

	claims := &accessTokenClaim{
		ID:      userID,
		Email:   email,
		KeyType: "access",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return tokenString, fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}

func (s *JWT) ValidateToken(tokenString string) (AccessTokenInfo, error) {
	accesTokenInfo := AccessTokenInfo{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&accessTokenClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(s.secretKey), nil
		})

	if err != nil {
		return accesTokenInfo, err
	}

	if !token.Valid {
		return accesTokenInfo, fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(*accessTokenClaim)
	if !ok {
		return accesTokenInfo, fmt.Errorf("couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return accesTokenInfo, fmt.Errorf("token expired")
	}

	accesTokenInfo.ID = claims.ID
	accesTokenInfo.Email = claims.Email

	return accesTokenInfo, nil
}

package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
	"errors"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	mySigningKey := []byte(tokenSecret)
	current_time := time.Now()
	expire_time := current_time.Add(expiresIn)
	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(current_time),
		ExpiresAt: jwt.NewNumericDate(expire_time),
		Issuer:    "chirpy",
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	return ss, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	} 
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		u, _ := uuid.Parse(claims.Subject)
		return u, nil
	} else {
		return uuid.Nil, errors.New("Unknown claims type, cannot proceed")
	}
	return uuid.Nil, nil
}
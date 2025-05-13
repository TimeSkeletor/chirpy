package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const signingKey = "JWTSigningKey"



func MakeJWT(userID uuid.UUID, tokenSecret string,  expiresIn time.Duration) (string, error) {
	issuedAt := time.Now().UTC()

	claims := jwt.RegisteredClaims{
		Issuer: 	"chirpy",
        IssuedAt:  	jwt.NewNumericDate(issuedAt),
		ExpiresAt: 	jwt.NewNumericDate(issuedAt.Add(expiresIn)),
		Subject: 	userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    signedToken, err := token.SignedString([]byte(tokenSecret))
    if err != nil {
        return "", err
    }

	return signedToken, nil
}


func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(tokenSecret), nil
		},
	)	
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("Invalid token")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("Couldn't parse claims")
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Couldn't parse user ID: %v", err)
	}

	return userId, nil
}
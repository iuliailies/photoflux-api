package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func parseJWT(tokenstr string, secret []byte) (Identity, bool) {

	// Parse the token to the standard Registered Claims.
	token, err := jwt.ParseWithClaims(tokenstr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})
	if err != nil {
		fmt.Printf("An unexpected error occured during authentication: %s.", err.Error())
		return Identity{}, false
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if ok && token.Valid {
		return Identity{
			User: uuid.MustParse(claims.Subject),
		}, true
	} else {
		return Identity{}, false
	}
}

func parseExpiredJWT(tokenstr string, secret []byte) (Identity, bool) {

	// Parse the token to the standard Registered Claims.
	token, err := jwt.ParseWithClaims(tokenstr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return secret, nil
	})

	if err != nil {
		var valerr *jwt.ValidationError
		// Works even if errors.As fails and valerr remains nil because of lazy evaluation.
		if errors.As(err, &valerr) && valerr.Errors == jwt.ValidationErrorExpired {
			// This is ok, so do not return the error.
			goto IF
		}
		fmt.Printf("parse expired: An unexpected error occured during authentication: %s.\n", err.Error())
		return Identity{}, false
	}
IF:
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		return Identity{
			User: uuid.MustParse(claims.Subject),
		}, true
	} else {
		return Identity{}, false
	}
}

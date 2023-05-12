package auth

import (
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

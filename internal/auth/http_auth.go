package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// An authentication middleware, decoding the authentication
// credentials and passing them along the gin context. It returns a 401
// Unauthorized response in case of invalid authentication data.
func BearerAuth(secret []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		tokenstr, err := getBearerToken(header)
		if err != nil {
			fmt.Println("error getting bearer", err)
			failAuth(ctx)
			return
		}
		identity, ok := parseJWT(tokenstr, secret)
		if !ok {
			fmt.Println("error parsing jwt", err)
			failAuth(ctx)
			return
		}
		ctx.Set(Authkey, identity)
	}

}

// failAuth fails an authentication request with a 401 Unauthorized.
func failAuth(ctx *gin.Context) {
	ctx.Header("WWW-Authenticate", "Bearer realm=\"Authorization required\"")
	// TODO: maybe add body
	ctx.AbortWithStatus(http.StatusUnauthorized)
}

// getBearerToken attempts to retrieve a token from a Bearer authentication request.
func getBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("received nil header")
	}

	content := strings.SplitN(header, " ", 2)
	if content[0] != "Bearer" {
		return "", fmt.Errorf("authentication method %s not supported", content[0])
	}

	return content[1], nil
}

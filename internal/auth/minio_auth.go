package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// MinioAuth is similar to BearerAuth middleware, except that it does not fail
// on invalid tokens since minio expects a specific answer format.
// Minio expects the token as a query parameter, not as a request header
func MinioAuth(secret []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenstr := ctx.Query("token")
		fmt.Println("tokenstr", tokenstr)
		identity, ok := parseJWT(tokenstr, secret)
		if !ok {
			fmt.Println("error parsing jwt")
			return
		}
		ctx.Set(Authkey, identity)
	}
}

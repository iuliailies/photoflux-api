package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iuliailies/photo-flux/internal/auth"
)

func (h *handler) HandleMinioAuth(ctx *gin.Context) {

	var identity auth.Identity
	var ok bool

	header, ok := ctx.Get(auth.Authkey)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"reason": "Token either invalid or not provided.",
		})
		return
	}

	identity, ok = header.(auth.Identity)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"reason": "Some internal server error occured while converting the token data.",
		})
		return
	}

	// Note that as an additional security constraint, one could check that
	// the user from the token indeed exists in the database. In this case
	// it is good enough, the token can't be changed by an attacker and if
	// it gets stolen it's game over anyway.

	ctx.JSON(http.StatusOK, gin.H{
		"user":               identity.User.String(),
		"maxValiditySeconds": h.minioTokenLifetime.Seconds(),
	})
}

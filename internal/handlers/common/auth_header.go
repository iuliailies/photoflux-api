package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iuliailies/photo-flux/internal/auth"
	"github.com/iuliailies/photo-flux/pkg/photoflux"
)

// GetAuthHeader receives the authentication header from the context and emits
// an error if it is missing or invalid. In that case the returned bool will be
// set to false and true if the header could be read correctly.
func GetAuthHeader(c *gin.Context) (auth.Identity, bool) {

	var ret auth.Identity

	header, ok := c.Get(auth.Authkey)
	if !ok {
		EmitError(c, photoflux.Error{
			Status: http.StatusUnauthorized,
			Title:  "Missing Auth Header",
			Detail: "The user identity had not been set.",
		})
		return ret, false
	}

	a, ok := header.(auth.Identity) // type assertion
	if !ok {
		EmitError(c, photoflux.Error{
			Status: http.StatusInternalServerError,
			Title:  "Invalid Auth Header",
			Detail: "Authentication header had unexpected type.",
		})
		return ret, false
	} else {
		return a, true
	}
}

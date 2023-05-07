package common

import (
	"github.com/gin-gonic/gin"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
)

// EmitError sends an error response back to the client.
func EmitError(ctx *gin.Context, err public.Error) {
	ctx.JSON(err.Status, public.ErrorResponse{
		Error: err,
	})
}

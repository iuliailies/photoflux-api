package handlers

import "github.com/gin-gonic/gin"

type PhotoHandler interface {
	HandleCreatePhoto(ctx *gin.Context)
	HandleUpdatePhoto(ctx *gin.Context)
	HandleListPhoto(ctx *gin.Context)
	HandleListMyPhoto(ctx *gin.Context)
	HandleGetPhoto(ctx *gin.Context)
	HandleDeletePhoto(ctx *gin.Context)
}

type UserHandler interface {
	HandleCreateUser(ctx *gin.Context)
}

type CategoryHandler interface {
	HandleListCategories(ctx *gin.Context)
}

type StarHandler interface {
	HandleStarPhoto(ctx *gin.Context)
	HandleIsPhotoStarred(ctx *gin.Context)
}

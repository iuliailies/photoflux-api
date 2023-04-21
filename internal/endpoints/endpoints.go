package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/iuliailies/photo-flux/internal/handlers"
)

func RegisterPhotos(router *gin.Engine, h handlers.PhotoHandler) {
	subrouter := router.Group("/api/photos")

	subrouter.GET("/", h.HandleListPhoto)
	subrouter.GET("/me/", h.HandleListMyPhoto)
	subrouter.GET("/:id", h.HandleGetPhoto)
	subrouter.POST("/", h.HandleCreatePhoto)
	subrouter.PATCH("/:id", h.HandleUpdatePhoto)
	subrouter.DELETE("/:id", h.HandleDeletePhoto)
}

func RegisterUsers(router *gin.Engine, h handlers.UserHandler) {
	subrouter := router.Group("/api/users")

	subrouter.POST("/", h.HandleCreateUser)
}

func RegisterCategories(router *gin.Engine, h handlers.CategoryHandler) {
	subrouter := router.Group("/api/categories")

	subrouter.GET("/", h.HandleListCategories)
}

func RegisterStars(router *gin.Engine, h handlers.StarHandler) {
	subrouter := router.Group("/api/stars")

	subrouter.POST("/", h.HandleStarPhoto)
	subrouter.GET("/", h.HandleIsPhotoStarred)
}

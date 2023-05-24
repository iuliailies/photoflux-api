package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/iuliailies/photo-flux/internal/auth"
	"github.com/iuliailies/photo-flux/internal/config"
	"github.com/iuliailies/photo-flux/internal/handlers"
)

func RegisterPhotos(router *gin.Engine, config config.Auth, h handlers.PhotoHandler) {
	subrouter := router.Group("/api/photos").Use(auth.BearerAuth(config.Secret))

	subrouter.GET("/", h.HandleListPhoto)
	subrouter.GET("/me/", h.HandleListMyPhoto)
	subrouter.GET("/:id", h.HandleGetPhoto)
	subrouter.POST("/", h.HandleCreatePhoto)
	subrouter.PATCH("/:id", h.HandleUpdatePhoto)
	subrouter.DELETE("/:id", h.HandleDeletePhoto)
}

func RegisterUsers(router *gin.Engine, config config.Auth, h handlers.UserHandler) {
	subrouter := router.Group("/api/users").Use(auth.BearerAuth(config.Secret))

	subrouter.GET("/:id", h.HandleGetUser)
}

func RegisterCategories(router *gin.Engine, config config.Auth, h handlers.CategoryHandler) {
	subrouter := router.Group("/api/categories").Use(auth.BearerAuth(config.Secret))

	subrouter.GET("/", h.HandleListCategories)
}

func RegisterStars(router *gin.Engine, config config.Auth, h handlers.StarHandler) {
	subrouter := router.Group("/api/stars").Use(auth.BearerAuth(config.Secret))

	subrouter.POST("/", h.HandleStarPhoto)
	subrouter.GET("/", h.HandleIsPhotoStarred)
}

func RegisterAuth(router *gin.Engine, config config.Auth, h handlers.AuthHandler) {
	subrouter := router.Group("/api/auth")
	subrouter.POST("/register", h.HandleRegister)
	subrouter.POST("/login", h.HandleLogin)
	subrouter.POST("/logout", h.HandleLogout)
	subrouter.POST("/minio", auth.MinioAuth(config.Secret), h.HandleMinioAuth)

	// the refresh request uses an independent middleware, which only allows expired tokens
	subsubrouter := subrouter.Group("/refresh").Use(auth.BearerAuthAllowExpired(config.Secret))
	subsubrouter.POST("", h.HandleRefresh)
}

func RegisterBoards(router *gin.Engine, config config.Auth, h handlers.BoardHandler) {
	subrouter := router.Group("/api/boards").Use(auth.BearerAuth(config.Secret))

	subrouter.GET("/", h.HandleListBoard)
	subrouter.GET("/:id", h.HandleGetBoard)
	subrouter.POST("/", h.HandleCreateBoard)
	subrouter.PATCH("/:id", h.HandleUpdateBoard)
}

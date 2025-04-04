package route

import (
	"api-gateway/docs"
	"api-gateway/internal/http/handlers"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoute(router *gin.Engine, authHandler *handlers.AuthHandler) {

	baseUrl := "/api/auth"

	docs.SwaggerInfo.BasePath = baseUrl
	docs.SwaggerInfo.Title = "Sso-omnia"

	auth := router.Group(baseUrl)

	auth.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	auth.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "healthy"})
	})

	auth.GET("/user", authHandler.GetUser)
	//auth.POST("/login", authHandler.Login)
	//auth.POST("/create", authHandler.CreateUser)
	auth.GET("/google/login", authHandler.AuthGoogleLogin)
	auth.GET("/google/callback", authHandler.AuthGoogleCallback)
	auth.GET("/vkid/login", authHandler.AuthVkidLogin)
	auth.GET("/vkid/callback", authHandler.AuthVkidCallback)
	//auth.POST("/change_password", authHandler.ChangePassword)
	auth.GET("/refresh", authHandler.Refresh)
	auth.GET("/validate", authHandler.Validate)
	auth.GET("/logout", authHandler.Logout)
}

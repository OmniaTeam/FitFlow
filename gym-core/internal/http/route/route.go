package route

import (
	"gym-core/docs"
	"gym-core/internal/http/handlers"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoute(router *gin.Engine, userHandler *handlers.UserHandler, exerciseHandler *handlers.ExerciseHandler, programHandler *handlers.ProgramHandler) {

	baseUrl := "/api/gym"

	docs.SwaggerInfo.BasePath = baseUrl
	docs.SwaggerInfo.Title = "Gym omnia"

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "healthy"})
	})

	// Группа маршрутов для пользователей
	users := router.Group("/users")
	{
		users.GET("/profile", userHandler.GetUserProfile)
		users.PUT("", userHandler.UpdateUser)
	}

	// Группа маршрутов для упражнений
	exercises := router.Group("/exercises")
	{
		exercises.GET("", exerciseHandler.GetAllExercises)
		exercises.GET("/:id", exerciseHandler.GetExerciseById)
	}

	// Группа маршрутов для программ тренировок
	programs := router.Group("/programs")
	{
		programs.GET("/user", programHandler.GetUserProgram)
	}
}

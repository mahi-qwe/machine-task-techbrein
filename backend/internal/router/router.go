package router

import (
	"net/http"

	"techbrein-project-management/internal/config"
	"techbrein-project-management/internal/domain/entities"
	"techbrein-project-management/internal/domain/services"
	"techbrein-project-management/internal/handlers"
	"techbrein-project-management/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(cfg config.Config, authService *services.AuthService, userService *services.UserService, projectService *services.ProjectService, taskService *services.TaskService) *gin.Engine {
	engine := gin.Default()

	engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", cfg.FrontendURL)
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	projectHandler := handlers.NewProjectHandler(projectService)
	taskHandler := handlers.NewTaskHandler(taskService)

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	engine.StaticFile("/swagger.json", "./docs/openapi.json")

	api := engine.Group("/api/v1")
	{
		api.POST("/auth/login", authHandler.Login)
	}

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.GET("/users", userHandler.List)
		protected.GET("/projects", projectHandler.List)
		protected.GET("/tasks", taskHandler.List)

		protected.PATCH("/tasks/:id/status", taskHandler.UpdateStatus)

		adminOnly := protected.Group("")
		adminOnly.Use(middleware.RequireRoles(string(entities.RoleAdmin)))
		{
			adminOnly.POST("/users", userHandler.Create)
			adminOnly.POST("/projects", projectHandler.Create)
			adminOnly.PUT("/projects/:id", projectHandler.Update)
			adminOnly.DELETE("/projects/:id", projectHandler.Delete)
			adminOnly.POST("/tasks", taskHandler.Create)
			adminOnly.PATCH("/tasks/:id/assign", taskHandler.Assign)
			adminOnly.PUT("/tasks/:id", taskHandler.Update)
		}
	}

	return engine
}

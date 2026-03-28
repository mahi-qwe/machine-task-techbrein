package main

import (
	"log"

	"techbrein-project-management/internal/config"
	"techbrein-project-management/internal/database"
	"techbrein-project-management/internal/domain/services"
	"techbrein-project-management/internal/repository/postgres"
	"techbrein-project-management/internal/router"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer func() {
		_ = database.Close(db)
	}()

	if err := database.Ping(db); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}

	if err := database.SeedAdmin(db, cfg); err != nil {
		log.Fatalf("failed to seed admin: %v", err)
	}

	userRepo := postgres.NewUserRepository(db)
	projectRepo := postgres.NewProjectRepository(db)
	taskRepo := postgres.NewTaskRepository(db)

	authService := services.NewAuthService(userRepo, cfg)
	userService := services.NewUserService(userRepo)
	projectService := services.NewProjectService(projectRepo)
	taskService := services.NewTaskService(db, taskRepo, projectRepo, userRepo)

	engine := router.Setup(cfg, authService, userService, projectService, taskService)

	log.Printf("server running on port %s", cfg.Port)
	if err := engine.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

package main

import (
	"log"

	"gin-template/internal/config"
	"gin-template/internal/database"
	"gin-template/internal/handler"
	"gin-template/internal/model"
	"gin-template/internal/repository"
	"gin-template/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.Connect(cfg.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(svc)

	r := handler.SetupRouter(userHandler)
	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

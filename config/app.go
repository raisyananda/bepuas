package config

import (
	"bepuas/app/repository"
	"bepuas/app/service"

	"database/sql"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppConfig struct {
	AuthService        *service.AuthService
	AchievementService *service.AchievementService
}

func NewConfig(pg *sql.DB, mongo *mongo.Database) *AppConfig {
	authRepo := repository.NewAuthRepository(pg)
	achievementRepo := repository.NewAchievementRepository(mongo)

	return &AppConfig{
		AuthService:        service.NewAuthService(authRepo),
		AchievementService: service.NewAchievementService(achievementRepo),
	}
}

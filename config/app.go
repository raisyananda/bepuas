package config

import (
	"bepuas/app/repository"
	"bepuas/app/service"

	"database/sql"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppConfig struct {
	AuthService             *service.AuthService
	AchievementMongoService *service.AchievementMongoService
	AchievementRefService   *service.AchievementRefService
	UserService             *service.UserService
	StudentService          *service.StudentService
	LectureService          *service.LecturerService
}

func NewConfig(pg *sql.DB, mongo *mongo.Database) *AppConfig {
	// REPOSITORY
	authRepo := repository.NewAuthRepository(pg)
	achievementMongoRepo := repository.NewAchievementMongoRepository(mongo)
	achievementRefRepo := repository.NewAchievementRefRepository(pg)
	userRepo := repository.NewUserRepository(pg)
	studentRepo := repository.NewStudentRepository(pg)
	lecturerRepo := repository.NewLecturerRepository(pg)

	// SERVICE
	authService := service.NewAuthService(authRepo)

	achievementMongoService := service.NewAchievementMongoService(
		achievementMongoRepo,
		achievementRefRepo,
		studentRepo,
	)

	achievementRefService := service.NewAchievementRefService(
		achievementRefRepo,
		achievementMongoService,
	)

	userService := service.NewUserService(userRepo)

	studentService := service.NewStudentService(
		studentRepo,
		achievementRefRepo)
	lecturerService := service.NewLecturerService(
		lecturerRepo,
		studentRepo)

	return &AppConfig{
		AuthService:             authService,
		AchievementMongoService: achievementMongoService,
		AchievementRefService:   achievementRefService,
		UserService:             userService,
		StudentService:          studentService,
		LectureService:          lecturerService,
	}
}

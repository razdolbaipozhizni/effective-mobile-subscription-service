package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/razdolbaipozhizni/effective-mobile-subscription-service/docs"
	"github.com/razdolbaipozhizni/effective-mobile-subscription-service/internal/config"
	"github.com/razdolbaipozhizni/effective-mobile-subscription-service/internal/handler"
	"github.com/razdolbaipozhizni/effective-mobile-subscription-service/internal/logger"
	"github.com/razdolbaipozhizni/effective-mobile-subscription-service/internal/model"
	"github.com/razdolbaipozhizni/effective-mobile-subscription-service/internal/repository"
	"github.com/razdolbaipozhizni/effective-mobile-subscription-service/internal/service"
)

// @title Subscription Service API
// @version 1.0
// @description REST API для агрегации данных об онлайн-подписках пользователей
// @host localhost:8080
// @BasePath /

func main() {
	// Инициализация конфига
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Инициализация логгера
	logger.Init(cfg.Logger.Level)
	log := logger.Get()

	log.Info().Msg("Starting Subscription Service...")

	// Подключение к PostgreSQL
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Автоматическая миграция
	err = db.AutoMigrate(&model.Subscription{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}

	log.Info().Msg("Database connected and migrated successfully")

	// Инициализация слоёв
	repo := repository.NewSubscriptionRepository(db)
	svc := service.NewSubscriptionService(repo)
	h := handler.NewSubscriptionHandler(svc)

	// Настройка Gin
	r := gin.Default()

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Роуты
	subscriptions := r.Group("/subscriptions")
	{
		subscriptions.POST("", h.Create)             // Create
		subscriptions.GET("", h.GetAll)             // List
		subscriptions.GET("/:id", h.GetByID)        // Read
		subscriptions.PUT("/:id", h.Update)         // Update
		subscriptions.DELETE("/:id", h.Delete)      // Delete
		subscriptions.GET("/total", h.GetTotalCost) // Специальный метод расчета
	}

	// Запуск сервера
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Info().Str("address", addr).Msg("Server starting")

	if err := r.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
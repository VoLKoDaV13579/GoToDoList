package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VoLKoDaV13579/GoToDoList/internal/config"
	"github.com/VoLKoDaV13579/GoToDoList/internal/database"
	"github.com/VoLKoDaV13579/GoToDoList/internal/handler"
	"github.com/VoLKoDaV13579/GoToDoList/internal/repository"
	"github.com/VoLKoDaV13579/GoToDoList/internal/service"
	"github.com/VoLKoDaV13579/GoToDoList/pkg/logger"
	"github.com/VoLKoDaV13579/GoToDoList/pkg/validator"
)

func main() {
	// Загружаем конфигурацию из переменных окружения
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Инициализируем логгер
	log := logger.NewTextLogger(cfg.App.Debug)
	log.Info("Starting application",
		"name", cfg.App.Name,
		"version", cfg.App.Version,
		"port", cfg.App.Port,
	)

	// Подключаемся к базе данных
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	log.Info("Database connection established")

	// Получаем generic database object для последующего закрытия соединения
	sqlDB, err := db.DB()
	if err != nil {
		log.Error("Failed to get database instance", "error", err)
		os.Exit(1)
	}

	// Инициализируем зависимости (Dependency Injection)
	// Репозиторий -> Сервис -> Handler (Clean Architecture)
	todoRepo := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)
	validatorInstance := validator.New()
	todoHandler := handler.NewTodoHandler(todoService, validatorInstance)

	// Настраиваем роутер с middleware
	router := handler.SetupRouter(todoHandler, log)

	// Создаем HTTP сервер
	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Info("Starting HTTP server", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	log.Info("Application started successfully")

	// Graceful Shutdown - ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	// Перехватываем SIGINT (Ctrl+C) и SIGTERM (docker stop, kubernetes)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Блокируемся до получения сигнала
	sig := <-quit
	log.Info("Received shutdown signal", "signal", sig.String())

	// Создаем контекст с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Останавливаем HTTP сервер с ожиданием завершения активных запросов
	log.Info("Shutting down HTTP server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	} else {
		log.Info("HTTP server stopped gracefully")
	}

	// Закрываем соединение с базой данных
	log.Info("Closing database connection...")
	if err := sqlDB.Close(); err != nil {
		log.Error("Error closing database connection", "error", err)
	} else {
		log.Info("Database connection closed")
	}

	log.Info("Application shutdown complete")
}

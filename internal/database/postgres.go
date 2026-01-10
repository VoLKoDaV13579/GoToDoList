package database

import (
	"fmt"
	"log"
	"time"

	"github.com/VoLKoDaV13579/GoToDoList/internal/config"
	"github.com/VoLKoDaV13579/GoToDoList/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresDB создает подключение к PostgreSQL и выполняет миграции
func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	// Формируем DSN (Data Source Name) для подключения
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	// Настраиваем уровень логирования GORM
	var gormLogger logger.Interface
	if cfg.App.Debug {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	// Открываем подключение к БД с настройками
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Получаем generic database object для настройки пула соединений
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Настраиваем пул соединений для production
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute)

	// Проверяем подключение
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL")

	// Выполняем автоматические миграции
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// runMigrations выполняет автоматические миграции моделей
func runMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Список моделей для миграции
	models := []interface{}{
		&model.Todo{},
	}

	// Выполняем AutoMigrate для каждой модели
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

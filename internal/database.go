package internal

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	migrate_postgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectToDB godoc
// @Summary Инициализация подключения к базе данных
// @Description Создает подключение к PostgreSQL используя DSN из переменных окружения
// @Tags database
// @Success 200 {object} map[string]interface{} "Успешное подключение"
// @Failure 500 {object} map[string]string "Ошибка подключения"
func ConnectToDB() (*gorm.DB, error) {
	dsn := (os.Getenv("DB_URL"))

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// RunMigrations godoc
// @Summary Применение миграций базы данных
// @Description Применяет все pending миграции из указанной директории
// @Tags database
// @Success 200 {object} map[string]interface{} "Миграции успешно применены"
// @Failure 500 {object} map[string]string "Ошибка миграции"
func RunMigrations(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	driver, err := migrate_postgres.WithInstance(sqlDB, &migrate_postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/internal/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	version, dirty, err := m.Version()
	log.Printf("Current migration version: %d, dirty: %v", version, dirty)

	return nil
}

// checkDBConnection godoc
// @Summary Проверка подключения к БД
// @Description Выполняет простой запрос для проверки работоспособности подключения
// @Tags database
// @Success 200 {object} map[string]string "Успешная проверка"
// @Failure 500 {object} map[string]string "Ошибка подключения"
func checkDBConnection(db *gorm.DB) error {
	var version string
	err := db.Raw("SELECT version()").Scan(&version).Error
	if err != nil {
		return err
	}
	log.Println("Connected to PostgreSQL:", version)
	return nil
}

package internal

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

func init() {
	// Инициализация базы данных и миграций
	dbOnce.Do(func() {
		var err error
		// Подключение к базе данных
		db, err = ConnectToDB()
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		// Проверка подключения
		if err := checkDBConnection(db); err != nil {
			log.Fatalf("Database connection check failed: %v", err)
		}

		// Выполнение миграций
		db.AutoMigrate(&Subscription{})
		if err := RunMigrations(db); err != nil {
			log.Fatalf("Migrations failed: %v", err)
		}
		log.Println("Database initialized and migrations applied successfully")
	})
}

// CreateSubscription godoc
// @Summary Create a new subscription
// @Description Create a new user subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body Subscription true "Subscription data"
// @Success 302 {string} string "Redirect to home page"
// @Failure 400 {object} map[string]string "Invalid request"
// @Router /create-subscription [post]
func CreateSubscription(c echo.Context) error {
	var Subscription Subscription
	err := c.Bind(&Subscription)
	if err != nil {
		log.Print("Ошибка привязки данных: ", err)
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}
	Subscription.CreatedAt, _ = Subscription.GetStartDate()
	Subscription.Ended, _ = Subscription.GetEndDate()
	Subscription.Ended = Subscription.Ended.AddDate(1, 0, 0)
	Subscription.EndTime = Subscription.Ended.Format("01-2006")

	result := db.Create(&Subscription)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	return c.Redirect(http.StatusFound, "/")
}

// ReadSubscription godoc
// @Summary Get subscription by ID
// @Description Get subscription details by ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} Subscription "Subscription details"
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 404 {object} map[string]string "Subscription not found"
// @Failure 500 {object} map[string]string "Database error"
// @Router /subscription/{id} [get]
func ReadSubscription(c echo.Context) error {
	id := c.Param("id")
	// Преобразование ID
	ID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Неверный ID подписки",
		})
	}
	var sub Subscription
	result := db.Where("id = ?", ID).First(&sub)
	if result.Error != nil {
		log.Printf("DB error: %v", result.Error)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Ошибка базы данных",
		})
	}

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Подписка не найдена",
		})
	}
	return c.Render(http.StatusOK, "subscription", map[string]any{
		"S": sub,
	})
}

// UpdateSubscription godoc
// @Summary Update subscription
// @Description Update existing subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param subscription body UpdateRequest true "Updated subscription data"
// @Success 200 {object} map[string]string "Update status"
// @Failure 400 {object} map[string]string "Invalid data"
// @Failure 404 {object} map[string]string "Subscription not found"
// @Failure 500 {object} map[string]string "Database error"
// @Router /update-subscription/{id} [put]
func UpdateSubscription(c echo.Context) error {
	id := c.Param("id")

	var req UpdateRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("Bind error: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Неверный формат данных: " + err.Error(),
		})
	}

	// Валидация
	if req.Price <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Цена должна быть положительной",
		})
	}
	// Преобразование ID
	subId, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Неверный ID подписки",
		})
	}

	req.CreatedAt, _ = req.GetStartDate()
	req.Ended, _ = req.GetStartDate()

	//  Обновление в БД
	result := db.Model(&Subscription{}).Where("id = ?", subId).Updates(map[string]interface{}{
		"service_name": req.ServiceName,
		"price":        req.Price,
		"user_id":      req.UserID,
		"start_date":   req.StartDate,
		"end_time":     req.EndTime,
		"created_at":   req.CreatedAt,
		"ended":        req.Ended,
	})

	if result.Error != nil {
		log.Printf("DB error: %v", result.Error)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Ошибка базы данных",
		})
	}

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Подписка не найдена",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}

// DeleteSubscription godoc
// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags subscriptions
// @Param id path int true "Subscription ID"
// @Success 204 "No content"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Subscription not found"
// @Failure 500 {object} map[string]string "Delete failed"
// @Router /delete-subscription/{id} [delete]
func DeleteSubscription(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	// Удаление из БД (пример с GORM)
	result := db.Delete(&Subscription{}, id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Delete failed"})
	}

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Subscription not found"})
	}

	return c.NoContent(http.StatusNoContent)
}

// List godoc
// @Summary List all subscriptions
// @Description Get list of all subscriptions
// @Tags subscriptions
// @Produce html
// @Success 200 {object} map[string]interface{} "List of subscriptions"
// @Failure 500 {object} map[string]string "Database error"
// @Router / [get]
func List(c echo.Context) error {
	// Получение всех записей
	var subscriptions []Subscription
	result := db.Find(&subscriptions)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	return c.Render(http.StatusOK, "list", map[string]any{
		"subs": subscriptions,
	})
}

// CalculateSubscriptionsSum godoc
// @Summary Calculate subscriptions total cost
// @Description Calculate total cost of subscriptions with filters
// @Tags analytics
// @Produce json
// @Param user_id query string false "User ID filter"
// @Param service_name query string false "Service name filter"
// @Param start_date query string true "Start date (format: MM-YYYY)"
// @Param end_date query string true "End date (format: MM-YYYY)"
// @Success 200 {object} map[string]interface{} "Total cost result"
// @Failure 400 {object} map[string]string "Invalid parameters"
// @Failure 500 {object} map[string]string "Calculation error"
// @Router /subscriptions/total [get]
func CalculateSubscriptionsSum(c echo.Context) error {
	// Парсим параметры запроса
	var req SubscriptionSumRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Неверные параметры запроса",
		})
	}

	// Валидация параметров
	if req.StartDate == "" || req.EndDate == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Необходимо указать период (start_date и end_date)",
		})
	}

	// Строим запрос к БД
	var query *gorm.DB
	query = db.Model(&Subscription{}).Select("COALESCE(SUM(price), 0) as total")

	req.BeginDate, _ = req.GetStartDate()
	req.EndedDate, _ = req.GetEndDate()
	// Добавляем фильтры, если они указаны
	if req.UserID != "" && req.ServiceName != "" {
		query = query.Where("user_id = ?", req.UserID).Where("created_at BETWEEN ? AND ? AND deleted_at IS NULL", req.BeginDate, req.EndedDate).Where("service_name = ?", req.ServiceName)
	} else if req.ServiceName != "" {
		query = query.Where("service_name = ?", req.ServiceName).Where("created_at BETWEEN ? AND ? AND deleted_at IS NULL", req.BeginDate, req.EndedDate)
	} else if req.UserID != "" {
		query = query.Where("user_id = ?", req.UserID).Where("created_at BETWEEN ? AND ? AND deleted_at IS NULL", req.BeginDate, req.EndedDate)
	} else {
		query = query.Where("created_at BETWEEN ? AND ? AND deleted_at IS NULL", req.BeginDate, req.EndedDate)
	}

	// Выполняем запрос
	var result struct {
		Total int `gorm:"column:total"`
	}
	if err := query.Scan(&result).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Ошибка при расчете суммы",
		})
	}

	// Возвращаем результат
	return c.JSON(http.StatusOK, map[string]interface{}{
		"total":    result.Total,
		"currency": "RUB",
		"period":   fmt.Sprintf("%s - %s", req.StartDate, req.EndDate),
		"user_id":  req.UserID,
		"service":  req.ServiceName,
	})
}

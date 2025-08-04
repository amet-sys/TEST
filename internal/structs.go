package internal

import (
	"time"

	"gorm.io/gorm"
)

// Subscription represents user subscription model
// @Description Модель подписки пользователя
// @ID Subscription
type Subscription struct {
	ID          uint   `json:"id" form:"id" gorm:"primaryKey"`
	ServiceName string `json:"service_name" form:"service_name" gorm:"not null"`
	Price       int    `json:"price" form:"price" gorm:"not null;check:price > 0"`
	User_ID     string `json:"user_id" form:"user_id" gorm:"type:uuid;not null"`
	StartDate   string `json:"start_date" form:"start_date" gorm:"not null"`
	EndTime     string
	Ended       time.Time
	CreatedAt   time.Time      `gorm:"not null"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	gorm.Model
}

// GetStartDate godoc
// @Summary Преобразование даты начала
// @Description Конвертирует строковое представление даты в формат time.Time
// @Tags internal
// @Param date body string true "Дата в формате MM-YYYY"
// @Success 200 {object} time.Time "Успешное преобразование"
// @Failure 400 {object} map[string]string "Неверный формат даты"
func (s *Subscription) GetStartDate() (time.Time, error) {
	return time.Parse("01-2006", s.StartDate)
}

// GetEndDate godoc
// @Summary Преобразование даты окончания
// @Description Конвертирует строковое представление даты в формат time.Time
// @Tags internal
// @Param date body string true "Дата в формате MM-YYYY"
// @Success 200 {object} time.Time "Успешное преобразование"
// @Failure 400 {object} map[string]string "Неверный формат даты"
func (s *Subscription) GetEndDate() (time.Time, error) {
	return time.Parse("01-2006", s.StartDate)
}

// UpdateRequest represents subscription update model
// @Description Модель для обновления подписки
type UpdateRequest struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndTime     string `json:"end_time"`
	Ended       time.Time
	CreatedAt   time.Time
}

func (s *UpdateRequest) GetStartDate() (time.Time, error) {
	return time.Parse("01-2006", s.StartDate)
}

func (s *UpdateRequest) GetEndDate() (time.Time, error) {
	return time.Parse("01-2006", s.EndTime)
}

// SubscriptionSumRequest represents subscription sum calculation request
// @Description Параметры запроса для расчета суммы подписок
type SubscriptionSumRequest struct {
	UserID      string `query:"user_id"`      // Фильтр по ID пользователя
	ServiceName string `query:"service_name"` // Фильтр по названию подписки
	StartDate   string `query:"start_date"`   // Начало периода
	EndDate     string `query:"end_date"`     // Конец периода
	BeginDate   time.Time
	EndedDate   time.Time
}

func (s *SubscriptionSumRequest) GetStartDate() (time.Time, error) {
	return time.Parse("01-2006", s.StartDate)
}

func (s *SubscriptionSumRequest) GetEndDate() (time.Time, error) {
	return time.Parse("01-2006", s.EndDate)
}

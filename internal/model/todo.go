package model

import (
	"time"

	"gorm.io/gorm"
)

// TodoStatus представляет статус задачи
type TodoStatus string

const (
	// StatusPending - задача ожидает выполнения
	StatusPending TodoStatus = "pending"
	// StatusInProgress - задача в процессе выполнения
	StatusInProgress TodoStatus = "in_progress"
	// StatusCompleted - задача завершена
	StatusCompleted TodoStatus = "completed"
)

// Todo представляет модель задачи в системе
type Todo struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Title       string         `gorm:"not null;size:255" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Status      TodoStatus     `gorm:"type:varchar(20);default:'pending';not null" json:"status"`
	Priority    int            `gorm:"default:0" json:"priority"` // 0-низкий, 1-средний, 2-высокий
	DueDate     *time.Time     `json:"due_date,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName возвращает имя таблицы для GORM
func (Todo) TableName() string {
	return "todos"
}

// CreateTodoRequest представляет запрос на создание задачи
type CreateTodoRequest struct {
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Description string     `json:"description" validate:"max=1000"`
	Status      TodoStatus `json:"status" validate:"omitempty,oneof=pending in_progress completed"`
	Priority    int        `json:"priority" validate:"min=0,max=2"`
	DueDate     *time.Time `json:"due_date"`
}

// UpdateTodoRequest представляет запрос на обновление задачи
type UpdateTodoRequest struct {
	Title       *string     `json:"title" validate:"omitempty,min=1,max=255"`
	Description *string     `json:"description" validate:"omitempty,max=1000"`
	Status      *TodoStatus `json:"status" validate:"omitempty,oneof=pending in_progress completed"`
	Priority    *int        `json:"priority" validate:"omitempty,min=0,max=2"`
	DueDate     *time.Time  `json:"due_date"`
}

// TodoFilter представляет фильтры для поиска задач
type TodoFilter struct {
	Status   *TodoStatus `form:"status"`
	Priority *int        `form:"priority"`
	Page     int         `form:"page" validate:"min=1"`
	PageSize int         `form:"page_size" validate:"min=1,max=100"`
}

// TodoListResponse представляет ответ со списком задач
type TodoListResponse struct {
	Data       []Todo `json:"data"`
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}

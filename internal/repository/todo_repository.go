package repository

import (
	"errors"
	"fmt"

	"github.com/VoLKoDaV13579/GoToDoList/internal/model"
	"gorm.io/gorm"
)

// TodoRepository определяет интерфейс для работы с задачами в БД
type TodoRepository interface {
	Create(todo *model.Todo) error
	GetByID(id uint) (*model.Todo, error)
	GetAll(filter *model.TodoFilter) ([]model.Todo, int64, error)
	Update(todo *model.Todo) error
	Delete(id uint) error
}

// todoRepository реализует интерфейс TodoRepository
type todoRepository struct {
	db *gorm.DB
}

// NewTodoRepository создает новый экземпляр репозитория задач
func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{
		db: db,
	}
}

// Create создает новую задачу в БД
func (r *todoRepository) Create(todo *model.Todo) error {
	if err := r.db.Create(todo).Error; err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}
	return nil
}

// GetByID получает задачу по ID
func (r *todoRepository) GetByID(id uint) (*model.Todo, error) {
	var todo model.Todo

	if err := r.db.First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("todo with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	return &todo, nil
}

// GetAll получает все задачи с учетом фильтров и пагинации
func (r *todoRepository) GetAll(filter *model.TodoFilter) ([]model.Todo, int64, error) {
	var todos []model.Todo
	var total int64

	// Базовый запрос
	query := r.db.Model(&model.Todo{})

	// Применяем фильтры
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Priority != nil {
		query = query.Where("priority = ?", *filter.Priority)
	}

	// Получаем общее количество записей (до пагинации)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count todos: %w", err)
	}

	// Применяем пагинацию
	offset := (filter.Page - 1) * filter.PageSize
	query = query.Offset(offset).Limit(filter.PageSize)

	// Сортировка: сначала по приоритету (DESC), потом по дате создания (DESC)
	query = query.Order("priority DESC, created_at DESC")

	// Выполняем запрос
	if err := query.Find(&todos).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get todos: %w", err)
	}

	return todos, total, nil
}

// Update обновляет существующую задачу
func (r *todoRepository) Update(todo *model.Todo) error {
	// Проверяем существование задачи
	var existing model.Todo
	if err := r.db.First(&existing, todo.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("todo with id %d not found", todo.ID)
		}
		return fmt.Errorf("failed to check todo existence: %w", err)
	}

	// Обновляем задачу
	if err := r.db.Save(todo).Error; err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}

	return nil
}

// Delete удаляет задачу (soft delete через GORM)
func (r *todoRepository) Delete(id uint) error {
	// Проверяем существование задачи
	var todo model.Todo
	if err := r.db.First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("todo with id %d not found", id)
		}
		return fmt.Errorf("failed to check todo existence: %w", err)
	}

	// Выполняем soft delete
	if err := r.db.Delete(&todo).Error; err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	return nil
}

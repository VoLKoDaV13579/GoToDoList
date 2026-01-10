package service

import (
	"fmt"
	"math"

	"github.com/VoLKoDaV13579/GoToDoList/internal/model"
	"github.com/VoLKoDaV13579/GoToDoList/internal/repository"
)

// TodoService определяет интерфейс бизнес-логики для работы с задачами
type TodoService interface {
	CreateTodo(req *model.CreateTodoRequest) (*model.Todo, error)
	GetTodoByID(id uint) (*model.Todo, error)
	GetAllTodos(filter *model.TodoFilter) (*model.TodoListResponse, error)
	UpdateTodo(id uint, req *model.UpdateTodoRequest) (*model.Todo, error)
	DeleteTodo(id uint) error
	UpdateTodoStatus(id uint, status model.TodoStatus) (*model.Todo, error)
}

// todoService реализует интерфейс TodoService
type todoService struct {
	repo repository.TodoRepository
}

// NewTodoService создает новый экземпляр сервиса задач
func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{
		repo: repo,
	}
}

// CreateTodo создает новую задачу с валидацией
func (s *todoService) CreateTodo(req *model.CreateTodoRequest) (*model.Todo, error) {
	// Устанавливаем значения по умолчанию, если не указаны
	status := req.Status
	if status == "" {
		status = model.StatusPending
	}

	// Создаем модель задачи
	todo := &model.Todo{
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}

	// Сохраняем в БД
	if err := s.repo.Create(todo); err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	return todo, nil
}

// GetTodoByID получает задачу по ID
func (s *todoService) GetTodoByID(id uint) (*model.Todo, error) {
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("todo not found: %w", err)
	}
	return todo, nil
}

// GetAllTodos получает список задач с пагинацией и фильтрами
func (s *todoService) GetAllTodos(filter *model.TodoFilter) (*model.TodoListResponse, error) {
	// Устанавливаем значения по умолчанию для пагинации
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 10
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100 // Ограничение максимального размера страницы
	}

	// Получаем данные из репозитория
	todos, total, err := s.repo.GetAll(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}

	// Рассчитываем общее количество страниц
	totalPages := int(math.Ceil(float64(total) / float64(filter.PageSize)))

	// Формируем ответ
	response := &model.TodoListResponse{
		Data:       todos,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}

	return response, nil
}

// UpdateTodo обновляет задачу частично (patch-запрос)
func (s *todoService) UpdateTodo(id uint, req *model.UpdateTodoRequest) (*model.Todo, error) {
	// Получаем существующую задачу
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("todo not found: %w", err)
	}

	// Обновляем только те поля, которые были переданы
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Status != nil {
		todo.Status = *req.Status
	}
	if req.Priority != nil {
		todo.Priority = *req.Priority
	}
	if req.DueDate != nil {
		todo.DueDate = req.DueDate
	}

	// Сохраняем изменения
	if err := s.repo.Update(todo); err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	return todo, nil
}

// DeleteTodo удаляет задачу
func (s *todoService) DeleteTodo(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}
	return nil
}

// UpdateTodoStatus обновляет только статус задачи (удобный метод)
func (s *todoService) UpdateTodoStatus(id uint, status model.TodoStatus) (*model.Todo, error) {
	// Получаем существующую задачу
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("todo not found: %w", err)
	}

	// Обновляем статус
	todo.Status = status

	// Сохраняем изменения
	if err := s.repo.Update(todo); err != nil {
		return nil, fmt.Errorf("failed to update todo status: %w", err)
	}

	return todo, nil
}

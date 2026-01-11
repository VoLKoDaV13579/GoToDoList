package repository

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VoLKoDaV13579/GoToDoList/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock, sqlDB
}

func TestNewTodoRepository(t *testing.T) {
	db, _, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewTodoRepository(db)
	assert.NotNil(t, repo)
}

func TestTodoRepository_Create(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewTodoRepository(db)

	t.Run("successful creation", func(t *testing.T) {
		now := time.Now()
		todo := &model.Todo{
			Title:       "Test Todo",
			Description: "Test Description",
			Status:      model.StatusPending,
			Priority:    1,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "todos"`)).
			WithArgs(
				todo.Title,
				todo.Description,
				todo.Status,
				todo.Priority,
				sqlmock.AnyArg(), // DueDate
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				sqlmock.AnyArg(), // DeletedAt
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := repo.Create(todo)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), todo.ID)
	})

	t.Run("creation error", func(t *testing.T) {
		todo := &model.Todo{
			Title: "Test Todo",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "todos"`)).
			WillReturnError(gorm.ErrInvalidData)
		mock.ExpectRollback()

		err := repo.Create(todo)
		assert.Error(t, err)
	})
}

func TestTodoRepository_GetByID(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewTodoRepository(db)

	t.Run("successful get", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "updated_at"}).
			AddRow(1, "Test Todo", "Test Description", model.StatusPending, 1, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE "todos"."id" = $1 AND "todos"."deleted_at" IS NULL ORDER BY "todos"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(rows)

		todo, err := repo.GetByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, todo)
		assert.Equal(t, uint(1), todo.ID)
		assert.Equal(t, "Test Todo", todo.Title)
	})

	t.Run("todo not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE "todos"."id" = $1 AND "todos"."deleted_at" IS NULL ORDER BY "todos"."id" LIMIT $2`)).
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		todo, err := repo.GetByID(999)
		assert.Error(t, err)
		assert.Nil(t, todo)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE "todos"."id" = $1 AND "todos"."deleted_at" IS NULL ORDER BY "todos"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnError(gorm.ErrInvalidDB)

		todo, err := repo.GetByID(1)
		assert.Error(t, err)
		assert.Nil(t, todo)
	})
}

func TestTodoRepository_GetAll(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewTodoRepository(db)

	t.Run("successful get all without filters", func(t *testing.T) {
		filter := &model.TodoFilter{
			Page:     1,
			PageSize: 10,
		}

		// Мок для Count
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "todos" WHERE "todos"."deleted_at" IS NULL`)).
			WillReturnRows(countRows)

		// Мок для Find
		rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "updated_at"}).
			AddRow(1, "Todo 1", "Desc 1", model.StatusPending, 2, time.Now(), time.Now()).
			AddRow(2, "Todo 2", "Desc 2", model.StatusInProgress, 1, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE "todos"."deleted_at" IS NULL ORDER BY priority DESC, created_at DESC LIMIT $1`)).
			WithArgs(10).
			WillReturnRows(rows)

		todos, total, err := repo.GetAll(filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, todos, 2)
	})

	t.Run("get all with status filter", func(t *testing.T) {
		status := model.StatusPending
		filter := &model.TodoFilter{
			Status:   &status,
			Page:     1,
			PageSize: 10,
		}

		// Мок для Count
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "todos" WHERE status = $1 AND "todos"."deleted_at" IS NULL`)).
			WithArgs(status).
			WillReturnRows(countRows)

		// Мок для Find
		rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "updated_at"}).
			AddRow(1, "Todo 1", "Desc 1", model.StatusPending, 1, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE status = $1 AND "todos"."deleted_at" IS NULL ORDER BY priority DESC, created_at DESC LIMIT $2`)).
			WithArgs(status, 10).
			WillReturnRows(rows)

		todos, total, err := repo.GetAll(filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, todos, 1)
	})

	t.Run("get all with pagination", func(t *testing.T) {
		filter := &model.TodoFilter{
			Page:     2,
			PageSize: 5,
		}

		// Мок для Count
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(10)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "todos" WHERE "todos"."deleted_at" IS NULL`)).
			WillReturnRows(countRows)

		// Мок для Find с offset
		rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "updated_at"})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE "todos"."deleted_at" IS NULL ORDER BY priority DESC, created_at DESC LIMIT $1 OFFSET $2`)).
			WithArgs(5, 5).
			WillReturnRows(rows)

		todos, total, err := repo.GetAll(filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(10), total)
		assert.Len(t, todos, 0)
	})

	t.Run("count error", func(t *testing.T) {
		filter := &model.TodoFilter{
			Page:     1,
			PageSize: 10,
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "todos" WHERE "todos"."deleted_at" IS NULL`)).
			WillReturnError(gorm.ErrInvalidDB)

		todos, total, err := repo.GetAll(filter)
		assert.Error(t, err)
		assert.Equal(t, int64(0), total)
		assert.Nil(t, todos)
	})
}

func TestTodoRepository_Update(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewTodoRepository(db)

	t.Run("successful update", func(t *testing.T) {
		todo := &model.Todo{
			ID:          1,
			Title:       "Updated Todo",
			Description: "Updated Description",
			Status:      model.StatusCompleted,
			Priority:    2,
		}

		// Мок для проверки существования
		rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "updated_at"}).
			AddRow(1, "Old Todo", "Old Desc", model.StatusPending, 1, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE "todos"."id" = $1 AND "todos"."deleted_at" IS NULL ORDER BY "todos"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(rows)

		// Мок для обновления
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "todos" SET`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(todo)
		assert.NoError(t, err)
	})

	t.Run("todo not found", func(t *testing.T) {
		todo := &model.Todo{
			ID:    999,
			Title: "Updated Todo",
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE "todos"."id" = $1 AND "todos"."deleted_at" IS NULL ORDER BY "todos"."id" LIMIT $2`)).
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		err := repo.Update(todo)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("update error", func(t *testing.T) {
		todo := &model.Todo{
			ID:    1,
			Title: "Updated Todo",
		}

		rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "updated_at"}).
			AddRow(1, "Old Todo", "Old Desc", model.StatusPending, 1, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE "todos"."id" = $1 AND "todos"."deleted_at" IS NULL ORDER BY "todos"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "todos" SET`)).
			WillReturnError(gorm.ErrInvalidData)
		mock.ExpectRollback()

		err := repo.Update(todo)
		assert.Error(t, err)
	})
}

func TestTodoRepository_Delete(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewTodoRepository(db)

	t.Run("successful delete", func(t *testing.T) {
		// Мок для проверки существования
		rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "updated_at"}).
			AddRow(1, "Todo to Delete", "Desc", model.StatusPending, 1, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE "todos"."id" = $1 AND "todos"."deleted_at" IS NULL ORDER BY "todos"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(rows)

		// Мок для soft delete
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "todos" SET "deleted_at"=$1 WHERE "todos"."id" = $2 AND "todos"."deleted_at" IS NULL`)).
			WithArgs(sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("todo not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE "todos"."id" = $1 AND "todos"."deleted_at" IS NULL ORDER BY "todos"."id" LIMIT $2`)).
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		err := repo.Delete(999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("delete error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "updated_at"}).
			AddRow(1, "Todo to Delete", "Desc", model.StatusPending, 1, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE "todos"."id" = $1 AND "todos"."deleted_at" IS NULL ORDER BY "todos"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "todos" SET "deleted_at"=$1 WHERE "todos"."id" = $2 AND "todos"."deleted_at" IS NULL`)).
			WillReturnError(gorm.ErrInvalidDB)
		mock.ExpectRollback()

		err := repo.Delete(1)
		assert.Error(t, err)
	})
}

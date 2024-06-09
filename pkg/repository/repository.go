package repository

import (
	"github.com/jmoiron/sqlx"
	"time"
	todoListSber "todo-list-sber"
)

type TodoItem interface {
	Create(item todoListSber.TodoItem) (int, error)
	GetAll() ([]todoListSber.TodoItem, error)
	GetById(id int) (todoListSber.TodoItem, error)
	Delete(id int) error
	Update(id int, input todoListSber.UpdateItemInput) error
	GetDoneTodoItems(date *time.Time, limit int, offset int) ([]todoListSber.TodoItem, error)
	GetUndoneTodoItems(date *time.Time, limit int, offset int) ([]todoListSber.TodoItem, error)
}

type Repository struct {
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		TodoItem: NewTodoItemPostgres(db),
	}
}

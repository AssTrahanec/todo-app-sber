package service

import (
	"time"
	"todo-list-sber/pkg/repository"
)
import todoListSber "todo-list-sber"

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type TodoItem interface {
	Create(todoItem todoListSber.TodoItem) (int, error)
	GetAll() ([]todoListSber.TodoItem, error)
	GetById(id int) (todoListSber.TodoItem, error)
	Delete(id int) error
	Update(id int, input todoListSber.UpdateItemInput) error
	GetDoneTodoItems(date *time.Time, limit int, offset int) ([]todoListSber.TodoItem, error)
	GetUndoneTodoItems(date *time.Time, limit int, offset int) ([]todoListSber.TodoItem, error)
}

type Service struct {
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		TodoItem: NewTodoItemService(repos.TodoItem),
	}
}

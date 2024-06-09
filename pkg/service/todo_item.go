package service

import (
	"time"
	todoListSber "todo-list-sber"
	"todo-list-sber/pkg/repository"
)

type TodoItemService struct {
	repo repository.TodoItem
}

func NewTodoItemService(repo repository.TodoItem) *TodoItemService {
	return &TodoItemService{repo: repo}
}
func (s *TodoItemService) Create(item todoListSber.TodoItem) (int, error) {
	return s.repo.Create(item)
}
func (s *TodoItemService) GetAll() ([]todoListSber.TodoItem, error) {
	return s.repo.GetAll()
}
func (s *TodoItemService) GetById(id int) (todoListSber.TodoItem, error) {
	return s.repo.GetById(id)
}
func (s *TodoItemService) Delete(id int) error {
	return s.repo.Delete(id)
}
func (s *TodoItemService) Update(id int, input todoListSber.UpdateItemInput) error {

	return s.repo.Update(id, input)
}
func (s *TodoItemService) GetDoneTodoItems(date *time.Time, limit int, offset int) ([]todoListSber.TodoItem, error) {
	return s.repo.GetDoneTodoItems(date, limit, offset)
}
func (s *TodoItemService) GetUndoneTodoItems(date *time.Time, limit int, offset int) ([]todoListSber.TodoItem, error) {
	return s.repo.GetUndoneTodoItems(date, limit, offset)
}

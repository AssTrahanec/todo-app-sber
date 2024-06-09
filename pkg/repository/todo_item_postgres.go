package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
	"time"
	todoListSber "todo-list-sber"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}
func (r *TodoItemPostgres) Create(item todoListSber.TodoItem) (int, error) {
	var id int
	createTodoItemQuery := "INSERT INTO todo_items (title, description, date, is_done) VALUES ($1, $2, $3, $4) RETURNING id;"
	err := r.db.QueryRowx(createTodoItemQuery, item.Title, item.Description, item.Date, item.IsDone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}
func (r *TodoItemPostgres) GetAll() ([]todoListSber.TodoItem, error) {
	var todoItems []todoListSber.TodoItem
	query := fmt.Sprintf("SELECT id, title, description, date, is_done FROM todo_items")
	err := r.db.Select(&todoItems, query)
	return todoItems, err
}
func (r *TodoItemPostgres) GetById(id int) (todoListSber.TodoItem, error) {

	var todoItem todoListSber.TodoItem
	query := fmt.Sprintf("SELECT id, title, description, date, is_done FROM todo_items where id = $1")
	err := r.db.Get(&todoItem, query, id)
	return todoItem, err
}
func (r *TodoItemPostgres) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM todo_items where id = $1")
	_, err := r.db.Exec(query, id)
	return err
}
func (r *TodoItemPostgres) Update(id int, input todoListSber.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}
	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}
	if input.IsDone != nil {
		setValues = append(setValues, fmt.Sprintf("is_done=$%d", argId))
		args = append(args, *input.IsDone)
		argId++
	}
	if input.Date != nil {
		setValues = append(setValues, fmt.Sprintf("date=$%d", argId))
		args = append(args, *input.Date)
		argId++
	}
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE todo_items SET %s WHERE id = $%d", setQuery, argId)
	args = append(args, id)
	log.Printf("SQL Query: %s", query)
	log.Printf("Arguments: %v", args)
	_, err := r.db.Exec(query, args...)
	return err
}
func (r *TodoItemPostgres) GetDoneTodoItems(date *time.Time, limit int, offset int) ([]todoListSber.TodoItem, error) {

	var todoItems []todoListSber.TodoItem
	query := "SELECT id, title, description, date, is_done FROM todo_items"
	args := []interface{}{}

	if date != nil {
		query += " WHERE date::date = $1 AND is_done = true"
		args = append(args, *date)
	} else {
		query += " WHERE is_done = $1"
		args = append(args, true)
	}
	query += " ORDER BY date OFFSET $2 LIMIT $3"
	args = append(args, offset, limit)

	err := r.db.Select(&todoItems, query, args...)
	return todoItems, err
}
func (r *TodoItemPostgres) GetUndoneTodoItems(date *time.Time, limit int, offset int) ([]todoListSber.TodoItem, error) {

	var todoItems []todoListSber.TodoItem
	query := "SELECT id, title, description, date, is_done FROM todo_items"
	args := []interface{}{}

	if date != nil {
		query += " WHERE date::date = $1 AND is_done = false"
		args = append(args, *date)
	} else {
		query += " WHERE is_done = $1"
		args = append(args, false)
	}
	query += " ORDER BY date OFFSET $2 LIMIT $3"
	args = append(args, offset, limit)

	err := r.db.Select(&todoItems, query, args...)
	return todoItems, err
}

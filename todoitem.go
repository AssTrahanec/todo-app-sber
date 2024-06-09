package todo_list_sber

import "time"

type TodoItem struct {
	Id          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Description string    `json:"description" db:"description"`
	Date        time.Time `json:"date" db:"date" binding:"required"`
	IsDone      bool      `json:"is_done" db:"is_done"`
}
type UpdateItemInput struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	IsDone      *bool      `json:"is_done"`
	Date        *time.Time `json:"date"`
}

package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	todoListSber "todo-list-sber"
)

type getAllTodoItemsResponse struct {
	Data []todoListSber.TodoItem `json:"data"`
}

// @Summary createTodoItem
// @Description create todo item date example: 2024-06-07T12:00:00Z
// @ID create-todo-item
// @Accept  json
// @Produce  json
// @Param input body todoListSber.TodoItem true "todo info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/todo [post]
func (h *Handler) createTodoItem(c *gin.Context) {
	var input todoListSber.TodoItem

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid input body")
		return
	}
	id, err := h.services.TodoItem.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// @Summary getAllTodoItems
// @Description get all todos
// @ID get-all-todo-items
// @Accept  json
// @Produce  json
// @Success 200 {array} todoListSber.TodoItem
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/todo [get]
func (h *Handler) getAllTodoItems(c *gin.Context) {
	todoItems, err := h.services.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": todoItems})
}

// @Summary getTodoItemById
// @Description get todo item by id
// @ID get-todo-item-by-id
// @Param id path string true "get todo by id"
// @Accept  json
// @Produce  json
// @Success 200 {object} todoListSber.TodoItem
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/todo/{id} [get]
func (h *Handler) getTodoItemById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	todoItem, err := h.services.GetById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": todoItem})
}

// @Summary updateTodoItem
// @Description update todo item
// @ID update-todo-item
// @Param id path string true "get todo by id"
// @Param input body todoListSber.UpdateItemInput true "todo info"
// @Accept  json
// @Produce  json
// @Success 200 {string} status ok
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/todo/{id} [put]
func (h *Handler) updateTodoItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	var input todoListSber.UpdateItemInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.IsDone == nil && input.Title == nil && input.Description == nil && input.Date == nil {
		newErrorResponse(c, http.StatusBadRequest, "EOF")
		return
	}
	err = h.services.Update(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// @Summary deleteTodoItem
// @Description delete todo item
// @ID delete-todo-item
// @Param id path string true "get todo by id"
// @Accept  json
// @Produce  json
// @Success 200 {string} status ok
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/todo/{id} [delete]
func (h *Handler) deleteTodoItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	err = h.services.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetDoneTodoItems
// @Tags get by is_done
// @Summary getDoneTodoItems
// @Description get done todos by date with pagination
// @ID get-done-todo-items
// @Accept  json
// @Produce  json
// @Param date query string false "Date in format YYYY-MM-DD"
// @Param limit query int true "Limit of items to return"
// @Param offset query int true "Offset of items to return"
// @Success 200 {array} todoListSber.TodoItem
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/todo/done [get]
func (h *Handler) GetDoneTodoItems(c *gin.Context) {
	dateStr := c.Query("date")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	var date *time.Time
	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		date = &parsedDate
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}

	todos, err := h.services.GetDoneTodoItems(date, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

// GetUndoneTodoItems
// @Tags get by is_done
// @Summary getUndoneTodoItems
// @Description get undone todos by date with pagination
// @ID get-undone-todo-items
// @Accept  json
// @Produce  json
// @Param date query string false "Date in format YYYY-MM-DD"
// @Param limit query int true "Limit of items to return"
// @Param offset query int true "Offset of items to return"
// @Success 200 {array} todoListSber.TodoItem
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/todo/undone [get]
func (h *Handler) GetUndoneTodoItems(c *gin.Context) {
	dateStr := c.Query("date")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	var date *time.Time
	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		date = &parsedDate
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}

	todos, err := h.services.GetUndoneTodoItems(date, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "todo-list-sber/docs"
	"todo-list-sber/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api")
	{
		todo := api.Group("/todo")
		{
			todo.POST("/", h.createTodoItem)
			todo.GET("/", h.getAllTodoItems)
			todo.GET("/:id", h.getTodoItemById)
			todo.DELETE("/:id", h.deleteTodoItem)
			todo.PUT("/:id", h.updateTodoItem)
			todo.GET("/done", h.GetDoneTodoItems)
			todo.GET("/undone", h.GetUndoneTodoItems)
		}
	}
	return router
}

package main

import (
	_ "github.com/lib/pq"
	"log"
	"time"
	todolistsber "todo-list-sber"
	_ "todo-list-sber/docs"
	"todo-list-sber/pkg/handler"
	"todo-list-sber/pkg/repository"
	"todo-list-sber/pkg/service"
)

// @title           Todo List API
// @version         1.0
// @description     This is a sample server Todo server.
// @host            localhost:8080
// @BasePath        /

func main() {
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     "todo-list-postgres",
		Port:     "5432",
		Username: "postgres",
		Password: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
	})
	for {
		if err := db.Ping(); err == nil {
			break
		}
		log.Println("Waiting for database to be ready...")
		time.Sleep(time.Second)
	}

	if err != nil {
		log.Fatalf("error initialiazing db: %s", err.Error())
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todolistsber.Server)
	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("Error starting server: %s", err.Error())
	}
}

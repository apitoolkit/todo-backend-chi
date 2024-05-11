package main

import (
	"os"

	"github.com/apitoolkit/todo-backend-chi/server"
	"github.com/apitoolkit/todo-backend-chi/server/controller"
	"github.com/apitoolkit/todo-backend-chi/server/database"
)

var (
	port    = "8080"
	baseURL = "http://localhost:" + port
)

func init() {
	if env := os.Getenv("PORT"); env != "" {
		port = env
	}
	if env := os.Getenv("BASE_URL"); env != "" {
		baseURL = env
	}
}

func main() {
	dbConfig := database.Config{
		BaseURL: baseURL,
	}
	db := database.New(dbConfig)
	c := controller.Todo(db)
	s := server.New(baseURL)


	s.SetupRoutes( c)
	s.Serve(port)
}

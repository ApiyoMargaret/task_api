package main

import (
	"log"
	"net/http"
	"task-api/handler"
	"task-api/middleware"
	"task-api/store"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	db := store.NewMemoryStore()
	authH := &handler.AuthHandler{Store: db}
	taskH := &handler.TaskHandler{Store: db}

	// Public Routes
	r.Post("/register", authH.Register)
	r.Post("/login", authH.Login)

	// Protected Routes
	r.Route("/tasks", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/", taskH.Create)
		r.Get("/", taskH.List)
		r.Get("/{id}", taskH.Get)
		r.Put("/{id}", taskH.Update)
		r.Delete("/{id}", taskH.Delete)
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
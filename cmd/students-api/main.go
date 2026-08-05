package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/roystondz/students-api/internal/config"
	"github.com/roystondz/students-api/internal/http/handlers/auth"
	"github.com/roystondz/students-api/internal/http/handlers/student"
	"github.com/roystondz/students-api/internal/storage/sqlite"
	"github.com/roystondz/students-api/internal/http/middleware"
)

func main() {
	cfg := config.MustLoad()

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal("failed to create storage: " + err.Error())
	}
	slog.Info("storage created successfully")

	router := http.NewServeMux()

	studentHandler := student.New(storage)
	authHandler := auth.New(storage, cfg.JWTSecret)
	authMiddleware := middleware.AuthMiddleware([]byte(cfg.JWTSecret))

	// Student Endpoints
	router.Handle("POST /api/students", authMiddleware(http.HandlerFunc(studentHandler.Create)))
	router.Handle("GET /api/students/{id}", authMiddleware(http.HandlerFunc(studentHandler.Get)))
	router.Handle("GET /api/students", authMiddleware(http.HandlerFunc(studentHandler.GetAll)))
	router.Handle("DELETE /api/students/{id}", authMiddleware(http.HandlerFunc(studentHandler.Delete)))
	router.Handle("PUT /api/students/{id}", authMiddleware(http.HandlerFunc(studentHandler.Update)))

	// Auth Endpoints
	router.Handle("POST /api/auth/signup", http.HandlerFunc(authHandler.SignUp))
	router.Handle("POST /api/auth/signin", http.HandlerFunc(authHandler.SignIn))
	router.Handle("POST /api/auth/logout", authMiddleware(http.HandlerFunc(authHandler.Logout)))

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.HTTPServer.Address, cfg.HTTPServer.Port),
		Handler: router,
	}

	done := make(chan os.Signal, 1)

	fmt.Println("Server is starting on: ", cfg.HTTPServer.Address, cfg.HTTPServer.Port)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	<-done

	fmt.Println("Shutting down server gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Error shutting down server: ", slog.Any("error", err))
	}

	slog.Info("Server stopped")
}

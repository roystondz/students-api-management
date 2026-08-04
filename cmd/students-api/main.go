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
)

func main() {
	cfg := config.MustLoad()

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal("failed to create storage: " + err.Error())
	}
	slog.Info("storage created successfully")

	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.Create(storage))
	router.HandleFunc("GET /api/students/{id}", student.Get(storage))
	router.HandleFunc("GET /api/students", student.GetAll(storage))
	router.HandleFunc("DELETE /api/students/{id}", student.Delete(storage))
	router.HandleFunc("PUT /api/students/{id}", student.Update(storage))

	router.HandleFunc("POST /api/auth/signup", auth.SignUp(storage))
	router.HandleFunc("POST /api/auth/signin", auth.SignIn(storage))

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

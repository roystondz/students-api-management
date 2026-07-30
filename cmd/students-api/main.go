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
	"github.com/roystondz/students-api/internal/http/handlers/student"
)

func main(){
	cfg := config.MustLoad()

	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.Create())


	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.HTTPServer.Address, cfg.HTTPServer.Port),
		Handler: router,
	}

	done :=make(chan os.Signal, 1)
	
	fmt.Println("Server is starting on: ", cfg.HTTPServer.Address, cfg.HTTPServer.Port)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, 
		syscall.SIGQUIT)

	go func (){
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	<- done

	fmt.Println("Shutting down server gracefully")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Error shutting down server: ", slog.Any("error", err))
	}

	slog.Info("Server stopped")
}
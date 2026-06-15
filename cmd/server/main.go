package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"ntuc/internal/config"
	"ntuc/internal/handler"
	"ntuc/internal/middleware"
	"ntuc/internal/repository"
	"ntuc/internal/service"
	"ntuc/pkg/logger"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	l, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer l.Sync()

	db, err := repository.NewMySQLDB(cfg.DatabaseURL)
	if err != nil {
		l.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	studentRepo := repository.NewStudentRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, l)
	studentService := service.NewStudentService(studentRepo, l)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, l)
	studentHandler := handler.NewStudentHandler(studentService, l)

	// Setup router
	router := http.NewServeMux()

	// Middleware stack
	var httpHandler http.Handler = router
	httpHandler = middleware.RequestID(httpHandler)
	httpHandler = middleware.Logging(httpHandler, l)
	httpHandler = middleware.Recovery(httpHandler, l)
	httpHandler = middleware.Timeout(httpHandler, 30*time.Second)
	httpHandler = middleware.CORS(httpHandler)

	// User Routes
	router.HandleFunc("GET /api/v1/users", userHandler.ListUsers)
	router.HandleFunc("GET /api/v1/users/{id}", userHandler.GetUser)
	router.HandleFunc("POST /api/v1/users", userHandler.CreateUser)
	router.HandleFunc("PUT /api/v1/users/{id}", userHandler.UpdateUser)
	router.HandleFunc("DELETE /api/v1/users/{id}", userHandler.DeleteUser)

	// Student Routes
	router.HandleFunc("GET /api/v1/students", studentHandler.ListStudents)
	router.HandleFunc("GET /api/v1/students/{id}", studentHandler.GetStudent)
	router.HandleFunc("POST /api/v1/students", studentHandler.CreateStudent)
	router.HandleFunc("PUT /api/v1/students/{id}", studentHandler.UpdateStudent)
	router.HandleFunc("DELETE /api/v1/students/{id}", studentHandler.DeleteStudent)

	// Health check
	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok"}`)
	})

	server := &http.Server{
		Addr:         cfg.ServerPort,
		Handler:      httpHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	l.Infof("server starting on %s", cfg.ServerPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		l.Fatalf("server error: %v", err)
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		l.Errorf("shutdown error: %v", err)
	}
}

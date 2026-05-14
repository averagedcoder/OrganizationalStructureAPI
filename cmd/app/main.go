package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/averagedcoder/OrganizationalStructureAPI/internal/config"
	"github.com/averagedcoder/OrganizationalStructureAPI/internal/database"
	"github.com/averagedcoder/OrganizationalStructureAPI/internal/handler"
	"github.com/averagedcoder/OrganizationalStructureAPI/internal/middleware"
	"github.com/averagedcoder/OrganizationalStructureAPI/internal/repository"
	"github.com/averagedcoder/OrganizationalStructureAPI/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	depRepo := repository.NewDepartmentRepository(db)
	depService := service.NewDepartmentService(depRepo)
	depHandler := handler.NewDepartmentHandler(depService)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/departments/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			depHandler.GetDepartment(w, r)
		case http.MethodPatch:
			depHandler.UpdateDepartment(w, r)
		case http.MethodDelete:
			depHandler.DeleteDepartment(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/departments", depHandler.CreateDepartment)

	loggedMux := middleware.Logger(mux)
	server := &http.Server{
		Addr:    ":8080",
		Handler: loggedMux,
	}

	go func() {
		log.Println("server started on :8080")

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("shutdown failed:", err)
	}

	log.Println("server stopped gracefully")
}

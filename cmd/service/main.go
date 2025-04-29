package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/hell-ecosystem/user-service/internal/config"
	httpdelivery "github.com/hell-ecosystem/user-service/internal/delivery/http"
	"github.com/hell-ecosystem/user-service/internal/repository/postgres"
	"github.com/hell-ecosystem/user-service/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DatabaseDSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := postgres.New(db)
	svc := service.NewUserService(repo)
	handler := httpdelivery.NewHandler(cfg, svc)

	server := &http.Server{
		Addr:         cfg.AppPort,
		Handler:      handler.Router(),
		ReadTimeout:  cfg.GetReadTimeout(),
		WriteTimeout: cfg.GetWriteTimeout(),
		IdleTimeout:  cfg.GetIdleTimeout(),
	}

	log.Printf("Starting %s on %s", cfg.ServiceName, cfg.AppPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

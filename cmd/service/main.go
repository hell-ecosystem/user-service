package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hell-ecosystem/user-service/internal/auth"
	"github.com/hell-ecosystem/user-service/internal/config"
	"github.com/hell-ecosystem/user-service/internal/delivery/httpdelivery"
	"github.com/hell-ecosystem/user-service/internal/repository/postgres"
	"github.com/hell-ecosystem/user-service/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db := mustInitPostgres(cfg)
	defer db.Close()

	authModule := auth.InitAuth(cfg)
	repo := postgres.New(db)
	svc := service.NewUserService(repo)
	handler := httpdelivery.NewHandler(svc, authModule)

	server := &http.Server{
		Addr:         cfg.AppPort,
		Handler:      handler.Router(),
		ReadTimeout:  cfg.GetReadTimeout(),
		WriteTimeout: cfg.GetWriteTimeout(),
		IdleTimeout:  cfg.GetIdleTimeout(),
	}

	// Обработка graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("Starting %s on %s", cfg.ServiceName, cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
}

func mustInitPostgres(cfg *config.Config) *sql.DB {
	db, err := sql.Open("postgres", cfg.DatabaseDSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(cfg.GetConnMaxLifetime())

	log.Printf("connected to postgres (maxConns=%d, idleConns=%d)\n",
		cfg.DBMaxOpenConns, cfg.DBMaxIdleConns)

	return db
}

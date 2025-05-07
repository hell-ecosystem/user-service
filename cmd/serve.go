package cmd

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hell-ecosystem/user-service/internal/config"
	"github.com/hell-ecosystem/user-service/internal/delivery/httpdelivery"
	"github.com/hell-ecosystem/user-service/internal/repository/postgres"
	"github.com/hell-ecosystem/user-service/internal/service"
	"github.com/spf13/cobra"

	"net/http"

	_ "github.com/lib/pq"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run service",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		// инициализация БД
		db := mustInitPostgres(cfg)
		defer db.Close()

		// репозиторий и сервис
		repo := postgres.New(db)
		svc := service.New(repo)

		handler := httpdelivery.NewHandler(svc)

		server := &http.Server{
			Addr:         cfg.AppPort,
			Handler:      handler.Router(),
			ReadTimeout:  cfg.GetReadTimeout(),
			WriteTimeout: cfg.GetWriteTimeout(),
			IdleTimeout:  cfg.GetIdleTimeout(),
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		go func() {
			log.Printf("Starting server on %s", cfg.AppPort)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("server error: %v", err)
			}
		}()

		<-ctx.Done()
		log.Println("Shutting down…")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("shutdown failed: %v", err)
		}
	},
}

func mustInitPostgres(cfg *config.Config) *sql.DB {
	dsn := cfg.DatabaseDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(cfg.GetConnMaxLifetime())
	log.Printf("connected to postgres (maxConns=%d, idleConns=%d)", cfg.DBMaxOpenConns, cfg.DBMaxIdleConns)
	return db
}

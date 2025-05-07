// cmd/serve.go
package cmd

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hell-ecosystem/user-service/internal/config"
	"github.com/hell-ecosystem/user-service/internal/db"
	"github.com/hell-ecosystem/user-service/internal/delivery/httpdelivery"
	"github.com/hell-ecosystem/user-service/internal/logger"
	"github.com/hell-ecosystem/user-service/internal/repository/postgres"
	"github.com/hell-ecosystem/user-service/internal/service"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start HTTP server",
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.Load()
		if err != nil {
			slog.Error("config load failed", slog.Any("error", err))
			os.Exit(1)
		}

		logger.InitLogger(cfg.GetLogLevel())
		slog.Info("logger initialized", slog.String("level", cfg.LogLevel))

		dbConn, err := db.Connect(cfg)
		if err != nil {
			slog.Error("db connect failed", slog.Any("error", err))
			os.Exit(1)
		}
		defer dbConn.Close()

		repo := postgres.New(dbConn)
		svc := service.New(repo)
		handler := httpdelivery.NewHandler(svc)

		server := &http.Server{
			Addr:         cfg.AppPort,
			Handler:      handler,
			ReadTimeout:  cfg.GetReadTimeout(),
			WriteTimeout: cfg.GetWriteTimeout(),
			IdleTimeout:  cfg.GetIdleTimeout(),
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		go func() {
			slog.Info("starting server", slog.String("addr", cfg.AppPort))
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				slog.Error("serve error", slog.Any("error", err))
				os.Exit(1)
			}
		}()

		<-ctx.Done()
		slog.Info("shutting down…")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("shutdown failed", slog.Any("error", err))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

// cmd/serve.go
package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hell-ecosystem/user-service/internal/config"
	"github.com/hell-ecosystem/user-service/internal/db"
	"github.com/hell-ecosystem/user-service/internal/delivery/httpdelivery"
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
			log.Fatalf("config load: %v", err)
		}

		// подключаемся к БД (с retry внутри)
		dbConn, err := db.Connect(cfg)
		if err != nil {
			log.Fatalf("db connect: %v", err)
		}
		defer dbConn.Close()

		repo := postgres.New(dbConn)
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
			log.Printf("serving on %s", cfg.AppPort)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("serve error: %v", err)
			}
		}()

		<-ctx.Done()
		log.Println("shutting down…")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("shutdown error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

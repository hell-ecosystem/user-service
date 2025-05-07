package cmd

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"

	"github.com/hell-ecosystem/user-service/internal/config"
	"github.com/hell-ecosystem/user-service/internal/db"
	"github.com/hell-ecosystem/user-service/internal/logger"
)

const migrationsDir = "migrations"

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "apply database migrations",
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("config load: %v", err)
		}

		logger.InitLogger(cfg.GetLogLevel())
		slog.Info("running migrations", slog.String("dir", migrationsDir))

		dbConn, err := db.Connect(cfg)
		if err != nil {
			log.Fatalf("db connect: %v", err)
		}
		defer dbConn.Close()

		goose.SetDialect("postgres")
		if err := goose.Up(dbConn, migrationsDir); err != nil {
			log.Fatalf("migrations failed: %v", err)
		}

		fmt.Println("migrations applied successfully")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

package cmd

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"

	"github.com/hell-ecosystem/user-service/internal/config"
)

const (
	// migrationsDir — жёстко зашитый путь к папке с миграциями
	migrationsDir = "migrations"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "make migrations",
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		db, err := sql.Open("postgres", cfg.DatabaseDSN())
		if err != nil {
			log.Fatalf("не удалось подключиться к БД: %v", err)
		}
		defer db.Close()

		goose.SetDialect("postgres")
		if err := goose.Up(db, migrationsDir); err != nil {
			log.Fatalf("goose up не удалось выполнить: %v", err)
		}

		fmt.Println("Миграции успешно применены")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

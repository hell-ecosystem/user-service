package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "user-service",
	Short: "User Service CLI",
	Long:  "CLI для запуска HTTP-сервера и выполнения миграций",
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(migrateCmd)
}

// Execute запускает корневую команду
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

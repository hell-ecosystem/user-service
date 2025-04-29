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
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := postgres.New(db)
	svc := service.NewUserService(repo)
	handler := httpdelivery.NewHandler(svc)

	server := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: handler.Router(),
	}

	log.Println("Starting server on", cfg.HTTPPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/http-server/handlers/auth"
	"awesomeProject/internal/http-server/handlers/ref"
	"awesomeProject/internal/http-server/handlers/reg"
	"awesomeProject/internal/http-server/services"
	"awesomeProject/internal/repo"
	"awesomeProject/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatal("Error loading .env file")
	}
	conf := config.LoadConfig()
	db, err := storage.NewPostgresDb(conf.Storage)
	if err != nil {
		logrus.Fatalf("failed to initialize db. %s", err.Error())
	}
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	repository := repo.NewRepository(db)
	service := services.NewService(repository)
	router.Get("/api/v1/auth", auth.NewAUTHHandler(*service))
	router.Post("/api/v1/reg", reg.NewRegHandler(*service))
	router.Post("/api/v1/refresh", ref.NewRefreshHandler(*service))
	logrus.Info("starting server")
	timeout, err := time.ParseDuration(conf.HttpServer.Timeout)
	if err != nil {

	}
	idleTimeout, err := time.ParseDuration(conf.HttpServer.IdleTimeout)
	if err != nil {

	}

	srv := &http.Server{
		Addr:              conf.HttpServer.Address,
		Handler:           router,
		ReadHeaderTimeout: timeout,
		WriteTimeout:      timeout,
		IdleTimeout:       idleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatalf("failed to start server")
	}
	logrus.Fatalf("server stopped")
}

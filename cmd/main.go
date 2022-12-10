package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bambi2/url-shorter/internal/config"
	"github.com/bambi2/url-shorter/internal/database"
	"github.com/bambi2/url-shorter/internal/handler"
	"github.com/bambi2/url-shorter/internal/repository"
	"github.com/bambi2/url-shorter/internal/server"
	"github.com/bambi2/url-shorter/internal/service"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	if err := config.InitConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s\n", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s\n", err.Error())
	}

	db, err := database.NewDatabase(os.Getenv("DB_TYPE"))
	if err != nil {
		logrus.Fatalf("error initializing to database: %s\n", err.Error())
	}

	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := handler.NewHandler(service)
	server := server.NewServer()

	go func() {
		if err := server.Run(viper.GetString("port"), handler.InitRoutes()); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("error occured while running http server: %s\n", err.Error())
		}
	}()

	logrus.Println("URL Shorter App Started...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logrus.Println("URL Shorter App Shutting Down...")

	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured while shutting down http server: %s\n", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured while closing database connection: %s\n", err.Error())
	}
}

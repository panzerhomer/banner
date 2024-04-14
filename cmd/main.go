package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/panzerhomer/banner/internal/config"
	"github.com/panzerhomer/banner/internal/handlers"
	repository "github.com/panzerhomer/banner/internal/repository/postgres"
	"github.com/panzerhomer/banner/internal/services"
)

var ctx = context.Background()

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("loading config failed")
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	if err = conn.Ping(ctx); err != nil {
		log.Fatalf("can't ping db: %s", err)
	}

	log.Println("database connected")

	bannerRepo := repository.NewBannerRepo(conn)
	bannerService := services.NewBannerService(bannerRepo)
	bannerHandler := handlers.NewBannerHandler(bannerService)
	routes := handlers.Routes(bannerHandler)

	httpServer := &http.Server{
		Addr:           cfg.Server.Address + ":" + cfg.Server.Port,
		Handler:        routes,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    time.Second * 5,
		WriteTimeout:   time.Second * 5,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	log.Println("server is running on " + cfg.Server.Address + ":" + cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("server is shutting down")

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("error occured on server shutting down: %s", err.Error())
	}
}

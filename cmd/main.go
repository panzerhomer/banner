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
	cache "github.com/panzerhomer/banner/internal/cache/redis.go"
	"github.com/panzerhomer/banner/internal/config"
	"github.com/panzerhomer/banner/internal/handlers"
	repository "github.com/panzerhomer/banner/internal/repository/postgres"
	"github.com/panzerhomer/banner/internal/services"
)

var ctx = context.Background()

const configPath = "./configs/config.yml"

func main() {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal("loading config failed: ", err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)

	log.Println("[init dsn]", dsn, "\n", cfg.Server)

	var conn *pgx.Conn

	maxAttempts := 10
	attempt := 1

	for {
		conn, err = pgx.Connect(ctx, dsn)
		if err != nil {
			log.Printf("attempt %d: unable to connect to database: %v\n", attempt, err)
			if attempt == maxAttempts {
				log.Fatalf("max attempts reached, unable to connect to database: %v\n", err)
			}
			attempt++
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	defer conn.Close(ctx)

	if err = conn.Ping(ctx); err != nil {
		log.Fatalf("can't ping db: %s", err)
	}

	log.Println("database connected")

	var redis *cache.Redis

	attempt = 1

	for {
		redis, err = cache.New(cfg)
		if err != nil {
			log.Printf("attempt %d: unable to connect to redis: %v\n", attempt, err)
			if attempt == maxAttempts {
				log.Fatalf("max attempts reached, unable to connect to redis: %v\n", err)
			}
			attempt++
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	defer redis.Disconnect()

	log.Println("redis connected")

	bannerRepo := repository.NewBannerRepo(conn)
	bannerService := services.NewBannerService(bannerRepo, redis)
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

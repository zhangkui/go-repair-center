package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	bootstrap "go-repair-center/internal/platform/bootstrap"
	"go-repair-center/internal/platform/config"
	"go-repair-center/internal/platform/database"
	"go-repair-center/internal/platform/logger"
	redisplatform "go-repair-center/internal/platform/redis"
	apihttp "go-repair-center/internal/transport/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	log := logger.New("debug", "json")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	mysqlDB, err := database.OpenMySQL(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer mysqlDB.Close()

	redisClient, err := redisplatform.Open(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer redisClient.Close()

	if err := database.RunMigrations(mysqlDB, "migrations"); err != nil {
		panic(err)
	}
	if err := bootstrap.SeedDefaults(ctx, mysqlDB, cfg); err != nil {
		panic(err)
	}

	handler := apihttp.NewRouter(cfg, mysqlDB, redisClient, log)
	srv := &http.Server{Addr: ":" + cfg.AppPort, Handler: handler, ReadTimeout: cfg.ReadTimeout, WriteTimeout: cfg.WriteTimeout}
	go func() { _ = srv.ListenAndServe() }()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
	log.WithField("addr", srv.Addr).Info("server stopped")
	fmt.Println("ok")
}

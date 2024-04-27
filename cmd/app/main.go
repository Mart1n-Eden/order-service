package main

import (
	"context"
	"log"
	"order-service/internal/cache"
	"order-service/internal/config"
	"order-service/internal/consumer"
	"order-service/internal/handler"
	"order-service/internal/server"
	"order-service/internal/service"
	"order-service/internal/storage"
	"os"
	"os/signal"
)

const (
	//shutdownTimeout = 5 * time.Second
	path = "./configs/config.yaml"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	pg, err := storage.NewPostgres(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	cach := cache.NewCache()
	stor := storage.NewStorage(pg)
	cons := consumer.NewConsumer(stor, cach)

	go func() {
		cons.StartSubscribe(cfg.Nats)
	}()

	srvc := service.NewService(cach)
	hand := handler.NewHandler(srvc)
	serv := server.NewServer(hand.InitRoutes(), cfg.Server.Port)

	go func() {
		log.Println("Running server")
		serv.Run()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit

	serv.Shutdown(ctx)
}

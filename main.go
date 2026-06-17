package main

import (
	"log"

	"interseguro/go-api/internal/application"
	"interseguro/go-api/internal/config"
	"interseguro/go-api/internal/infrastructure"
	httptransport "interseguro/go-api/internal/transport/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	statsClient := infrastructure.NewNodeStatsClient(cfg.NodeAPIURL)
	matrixService := application.NewMatrixService(statsClient)
	app := httptransport.NewServer(httptransport.Config{
		JWTSecret: cfg.JWTSecret,
	}, matrixService)

	log.Fatal(app.Listen(":" + cfg.Port))
}

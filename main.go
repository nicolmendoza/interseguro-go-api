package main

import (
	"log"

	matrixapp "interseguro/go-api/internal/application/matrix"
	"interseguro/go-api/internal/infrastructure/config"
	"interseguro/go-api/internal/infrastructure/statsclient"
	httptransport "interseguro/go-api/internal/transport/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	statsClient := statsclient.NewNodeStatsClient(cfg.NodeAPIURL)
	matrixService := matrixapp.NewService(statsClient)
	app := httptransport.NewServer(httptransport.Config{
		JWTSecret: cfg.JWTSecret,
	}, matrixService)

	log.Fatal(app.Listen(":" + cfg.Port))
}

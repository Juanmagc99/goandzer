package main

import (
	"fmt"
	"juanmagc99/goandzer/internal/balancer"
	"juanmagc99/goandzer/internal/proxy"
	"juanmagc99/goandzer/internal/registry"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg, err := registry.LoadConfig("config/config.yaml")
	fmt.Println(cfg)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	balancers := make(map[string]*balancer.RoundRobin)
	for _, svc := range cfg.Services {
		balancers[svc.PathPrefix] = balancer.NewRoundRobin(svc.Targets)
	}

	e.Use(proxy.NewReverseProxyMiddleware(balancers))

	if err := e.Start(cfg.BindAddress); err != nil {
		log.Fatal(err)
	}
}

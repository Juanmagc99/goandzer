package main

import (
	"fmt"
	"juanmagc99/goandzer/internal/auth"
	"juanmagc99/goandzer/internal/balancer"
	"juanmagc99/goandzer/internal/proxy"
	"juanmagc99/goandzer/internal/registry"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	adminKey := os.Getenv("ADMIN_API_KEY")
	if adminKey == "" {
		log.Fatal("missing ADMIN_API_KEY")
	}

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
		rr := balancer.NewRoundRobin(svc.Targets)
		balancers[svc.PathPrefix] = rr
		go registry.StartHealthChecker(svc, cfg.HealthCheck.Path, cfg.HealthCheck.Interval, rr)
	}

	e.Use(proxy.NewReverseProxyMiddleware(balancers))

	admin := e.Group("/admin", auth.APIKeyAuthMiddleware(adminKey))

	admin.GET("/status", func(c echo.Context) error {
		status := make(map[string][]string)
		for prefix, rr := range balancers {
			status[prefix] = rr.GetTargets()
		}
		return c.JSON(http.StatusOK, echo.Map{
			"services": status,
		})
	})

	if err := e.Start(cfg.BindAddress); err != nil {
		log.Fatal(err)
	}
}

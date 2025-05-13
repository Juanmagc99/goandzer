package main

import (
	"fmt"
	"juanmagc99/goandzer/internal/balancer"
	"juanmagc99/goandzer/internal/registry"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

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
	log.Println(balancers)

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path
			log.Printf(">> peticion entrante: %s", path)
			for prefix, rr := range balancers {
				log.Printf("   – probando prefijo %q …", prefix)
				if strings.HasPrefix(path, prefix) {
					target := rr.Next()
					if target == "" {
						return echo.NewHTTPError(http.StatusServiceUnavailable, "No servers available for this service")
					}

					u, err := url.Parse(target)
					if err != nil {
						return echo.NewHTTPError(http.StatusBadGateway, "Invalid target url")
					}

					proxy := httputil.NewSingleHostReverseProxy(u)

					originalDirector := proxy.Director
					proxy.Director = func(req *http.Request) {
						originalDirector(req)
						req.Header.Set("X-Forwarded-Host", req.Host)
						req.Header.Set("X-Origin-Host", u.Host)
					}

					proxy.ServeHTTP(c.Response(), c.Request())
					return nil
				}
			}

			return next(c)
		}
	})

	if err := e.Start(cfg.BindAddress); err != nil {
		log.Fatal(err)
	}
}

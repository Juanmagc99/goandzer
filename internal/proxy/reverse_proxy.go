package proxy

import (
	"juanmagc99/goandzer/internal/balancer"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

func NewReverseProxyMiddleware(balancers map[string]*balancer.RoundRobin) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path

			for prefix, rr := range balancers {
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
	}
}

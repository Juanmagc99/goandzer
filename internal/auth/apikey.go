package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func APIKeyAuthMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := c.Request().Header.Get("X-Api-Key")
			if key == "" || key != secret {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or missing API key")
			}

			return next(c)
		}
	}
}

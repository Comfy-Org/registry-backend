package drip_middleware

import (
	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
)

func ErrorLoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			if err != nil {
				log.Ctx(c.Request().Context()).
					Error().
					Err(err).
					Msgf("Error occurred Path: %s, Method: %s\n", c.Path(), c.Request().Method)
			}

			return err
		}
	}
}

package middleware

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"io"
)

func RequestLoggerMiddleware() echo.MiddlewareFunc {
	return echo_middleware.RequestLoggerWithConfig(echo_middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v echo_middleware.RequestLoggerValues) error {
			// Read the request body for logging
			requestBody, err := io.ReadAll(c.Request().Body)
			if err != nil {
				log.Ctx(c.Request().Context()).Error().Err(err).Msg("Failed to read request body")
				return err
			}
			// Reset the body for further use
			c.Request().Body = io.NopCloser(bytes.NewReader(requestBody))

			// Log request details including query parameters
			log.Ctx(c.Request().Context()).
				Info().
				Str("Method", c.Request().Method).
				Str("Path", c.Path()).
				Str("QueryParams", fmt.Sprintf("%v", c.QueryParams())).
				Str("RequestBody", string(requestBody)).
				Str("Headers", fmt.Sprintf("%v", c.Request().Header)).
				Msg("Request received")
			return nil
		},
	})
}

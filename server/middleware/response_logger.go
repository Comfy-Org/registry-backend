package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"registry-backend/tracing"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// Custom response writer to capture response body
type responseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (rw *responseWriter) Write(p []byte) (n int, err error) {
	// Capture the response body in the buffer
	n, err = rw.body.Write(p)
	if err != nil {
		return n, err
	}
	// Write to the actual ResponseWriter
	return rw.ResponseWriter.Write(p)
}

// ResponseLoggerMiddleware will log response details and errors.
func ResponseLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer tracing.TraceDefaultSegment(c.Request().Context(), "DripStrictServerImplementation.ResponseLoggerMiddleware")()

			// Create a custom response writer to capture the response body
			rw := &responseWriter{
				ResponseWriter: c.Response().Writer,
				body:           new(bytes.Buffer),
			}
			c.Response().Writer = rw

			// Call the next handler in the chain
			err := next(c)

			// Log any errors that occur during handling
			if err != nil {
				log.Ctx(c.Request().Context()).
					Error().
					Err(err).
					Str("Method", c.Request().Method).
					Str("Path", c.Path()).
					Msg("Error occurred during request handling")
			}

			// Log the response details
			log.Ctx(c.Request().Context()).
				Info().
				Int("Status", c.Response().Status).
				Str("ResponseBody", rw.body.String()).
				Str("ResponseHeaders", fmt.Sprintf("%v", c.Response().Header())).
				Msg("Response sent")

			return err
		}
	}
}

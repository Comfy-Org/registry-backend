package middleware

import (
	"context"
	drip_logging "registry-backend/logging"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

const CorrelationIDKey = "x-correlation-id"
const EndpointKey = "endpoint"

func generateFallbackCorrelationID() string {
	return "fallback-" + uuid.New().String()
}

// TracingMiddleware is a middleware that adds a trace ID to the context
func TracingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		traceID := c.Request().Header.Get("X-Cloud-Trace-Context")

		if traceID == "" {
			// Generate a fallback ID if no trace ID is present
			traceID = generateFallbackCorrelationID()
		}

		// Set trace and span IDs in the context
		reqCtx := context.WithValue(c.Request().Context(), CorrelationIDKey, traceID)

		// Set the endpoint in the context
		endpoint := c.Path()
		reqCtx = context.WithValue(reqCtx, EndpointKey, endpoint)

		requestLogger := drip_logging.SetupLogger()
		requestLogger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.
				Str(CorrelationIDKey, traceID).
				Str(EndpointKey, endpoint)
		})

		c.SetRequest(c.Request().WithContext(requestLogger.WithContext(reqCtx)))
		return next(c)
	}
}

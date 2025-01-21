package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
)

type teeReader struct {
	buf bytes.Buffer
	io.ReadCloser
}

func (tr *teeReader) Read(p []byte) (n int, err error) {
	n, err = tr.ReadCloser.Read(p)
	if err != nil || n == 0 {
		return
	}

	tr.buf.Write(p[:n])
	return
}

func (tr *teeReader) bytes() ([]byte, error) {
	if tr.buf.Len() == 0 {
		return io.ReadAll(tr.ReadCloser)
	}
	return tr.buf.Bytes(), nil
}

func RequestLoggerMiddleware() echo.MiddlewareFunc {
	ctxKey := struct{}{}

	rlw := echo_middleware.RequestLoggerWithConfig(echo_middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v echo_middleware.RequestLoggerValues) error {
			if txn := newrelic.FromContext(c.Request().Context()); txn != nil {
				segment := txn.StartSegment("RequestLoggerMiddleware")
				defer segment.End()
			}
			// Get the teeReader
			reader, ok := c.Request().Context().Value(ctxKey).(*teeReader)
			if !ok {
				reader = &teeReader{ReadCloser: c.Request().Body}
			}

			// Read the request body for logging
			requestBody, err := reader.bytes()
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

	mw := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if txn := newrelic.FromContext(c.Request().Context()); txn != nil {
				segment := txn.StartSegment("RequestLoggerMiddleware")
				defer segment.End()
			}

			req := c.Request()
			reader := &teeReader{ReadCloser: req.Body}
			req.Body = io.NopCloser(reader)
			c.SetRequest(req.WithContext(context.WithValue(req.Context(), ctxKey, reader)))
			return rlw(next)(c)
		}
	}

	return mw
}

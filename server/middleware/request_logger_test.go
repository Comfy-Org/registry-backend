package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tArrayWriter struct {
	logs [][]byte
}

func (w *tArrayWriter) Write(p []byte) (n int, err error) {
	w.logs = append(w.logs, p)
	return len(p), nil
}

func TestRequestLogger(t *testing.T) {
	jsonBody, err := json.Marshal(map[string]interface{}{
		"status": "status",
		"name":   "name",
	})
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/publishers", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	e := echo.New()

	e.Use(RequestLoggerMiddleware())

	t.Run("BodyHasBeenRead", func(t *testing.T) {
		writer := &tArrayWriter{}
		logger := zerolog.New(writer).With().Timestamp().Logger()
		ctx := logger.WithContext(context.Background())
		req := req.WithContext(ctx)

		var readBody []byte
		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) (err error) {
				readBody, err = io.ReadAll(c.Request().Body)
				if err != nil {
					return err
				}
				return fmt.Errorf("terminate")
			}
		})
		e.ServeHTTP(res, req)

		assert.Equal(t, jsonBody, readBody, "the body should still be readable after RequestLoggerMiddleware")
		require.Len(t, writer.logs, 1)
		logged := map[string]interface{}{}
		require.NoError(t, json.Unmarshal(writer.logs[0], &logged))
		require.Contains(t, logged, "RequestBody")
		assert.Contains(t, logged["RequestBody"].(string), string(jsonBody), "should contain the request body")
	})

	t.Run("BodyHasNotBeenRead", func(t *testing.T) {
		writer := &tArrayWriter{}
		logger := zerolog.New(writer).With().Timestamp().Logger()
		ctx := logger.WithContext(context.Background())
		req := req.WithContext(ctx)

		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) (err error) {
				return fmt.Errorf("terminate")
			}
		})
		e.ServeHTTP(res, req)

		require.Len(t, writer.logs, 1)
		logged := map[string]interface{}{}
		require.NoError(t, json.Unmarshal(writer.logs[0], &logged))
		require.Contains(t, logged, "RequestBody")
		assert.NotContains(t, logged["RequestBody"].(string), string(jsonBody), "should contain the request body")
	})
}

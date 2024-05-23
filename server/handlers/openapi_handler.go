package handler

import (
	"net/http"
	"registry-backend/drip"

	"github.com/labstack/echo/v4"
	// other imports
)

func SwaggerHandler(c echo.Context) error {
	swagger, err := drip.GetSwagger()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, swagger)
}

package web

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
)

func healthcheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

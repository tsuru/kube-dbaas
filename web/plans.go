package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tsuru/kube-dbaas/types"
)

func servicePlans(c echo.Context) error {
	return c.JSON(http.StatusOK, []types.Plan{
		{
			Name:        "mongodb-4-2-c0.1m0.2",
			Description: "Mongo 4.2 - 10% single CPU, 256Mi RAM",
		},

		{
			Name:        "mongodb-4-2-c1m2",
			Description: "Mongo 4.2 - 100% single CPU, 2GB RAM",
		},
	})
}

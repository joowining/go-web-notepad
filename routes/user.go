package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func UserRouterGroup(e *echo.Echo) {
	routeGroup := e.Group("/user")

	routeGroup.GET("", getUserName)
}

func getUserName(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"user": "kim"})
}
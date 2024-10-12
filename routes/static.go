package routes

import (
	"github.com/labstack/echo/v4"
)

func RegisterStatic(e *echo.Echo) {
	e.Static("/", "../html/index.html")
	e.File("/", "../html/index.html")
}

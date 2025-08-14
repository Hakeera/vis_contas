// Package routes
package routes

import (
	"vis_contas/internal/controller"

	"github.com/labstack/echo/v4"
)



func SetUpRoutes(e *echo.Echo)  {
	
	e.GET("/", controller.Home)
}

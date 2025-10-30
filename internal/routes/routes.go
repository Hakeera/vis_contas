// Package routes
package routes

import (
	"vis_contas/internal/auth"
	"vis_contas/internal/controller"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetUpRoutes(e *echo.Echo) {

	// Middleware Global Echo
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Público
	e.GET("/user", controller.LoginPage)
	e.POST("/login", controller.Login)

	// Protegido com JWT
	r := e.Group("")
	r.Use(auth.RequireAuth)

	r.GET("/", controller.Home)
	r.GET("/load_table", controller.LoadTable)
	r.PUT("/alternar_sit/:id", controller.AlternarSituacao)
	r.PUT("/alternar_sit/:id", controller.AlternarSituacao)
	r.POST("/logout", controller.Logout)
}

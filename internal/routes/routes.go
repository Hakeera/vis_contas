// Package routes
package routes

import (
	"vis_contas/internal/controller"

	"github.com/labstack/echo/v4"
)



func SetUpRoutes(e *echo.Echo)  {

	// Página Inicial
	e.GET("/", controller.Home)

	// Carrega a tabela com dados filtrados 
	e.GET("/load_table", controller.LoadTable)
	// Altera situação da fatura 
	e.PUT("/alterar_situacao/:id", controller.AlterarSituacao)
}

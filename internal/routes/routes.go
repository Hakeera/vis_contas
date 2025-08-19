// Package routes
package routes

import (
	"vis_contas/internal/controller"

	"github.com/labstack/echo/v4"
)



func SetUpRoutes(e *echo.Echo)  {

	// Página Inicial
	e.GET("/", controller.Home)

	// Filtros
	e.GET("/filtro_categoria", controller.FiltroCategoria)
	e.GET("/filtro_situacao", controller.FiltroSitucao)
	e.GET("/filtro_data", controller.FiltroData)
	e.GET("/limpar_filtros", controller.LimparFiltros)
}

// Package controller
package controller

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
)

// TODO Dinâmica de filtros em conjunto. Atualmente cada filtro atua de forma independente dos valores
// dos demais filtros. Requisições devem encaminhar todos os valores de todos os filtros


// LimparFiltros retorna os dados sem nenhum filtro 
func LimparFiltros(c echo.Context) error {

	tmpl, err := template .ParseFiles("view/invoice_table.html",
					"view/template/base.html",
)
	if err != nil {
		log.Println("Erro ao carregar templates:", err)
		return c.String(http.StatusInternalServerError, "Erro ao carregar templates: "+err.Error())
	}

	return tmpl.ExecuteTemplate(c.Response(), "invoice_table", nil)
}

// FiltroCategoria retorna os dados de acordo com a categoria selecionada
func FiltroCategoria(c echo.Context) error {
	
	categoria := c.QueryParam("categoria")
	fmt.Println("Categoria: ", categoria)

	tmpl, err := template .ParseFiles("view/invoice_table.html",
					"view/template/base.html",
)
	if err != nil {
		log.Println("Erro ao carregar templates:", err)
		return c.String(http.StatusInternalServerError, "Erro ao carregar templates: "+err.Error())
	}

	data := map[string]any{
		"Categoria": categoria,
	}

	return tmpl.ExecuteTemplate(c.Response(), "invoice_table", data)
}

// FiltroSitucao retorna os dados de acordo com a situação selecionada
func FiltroSitucao(c echo.Context) error { 
	
	situacao := c.QueryParam("situacao")
	fmt.Println("Situação: ", situacao)

	tmpl, err := template .ParseFiles("view/invoice_table.html",
					"view/template/base.html",
)
	if err != nil {
		log.Println("Erro ao carregar templates:", err)
		return c.String(http.StatusInternalServerError, "Erro ao carregar templates: "+err.Error())
	}

	data := map[string]any{
		"Situacao": situacao,
	}

	return tmpl.ExecuteTemplate(c.Response(), "invoice_table", data)
}

// FiltroData retorna os dados de acordo com a data selecionada
func FiltroData(c echo.Context) error {

	tmpl, err := template .ParseFiles("view/invoice_table.html",
					"view/template/base.html",
)
	if err != nil {
		log.Println("Erro ao carregar templates:", err)
		return c.String(http.StatusInternalServerError, "Erro ao carregar templates: "+err.Error())
	}

	return tmpl.ExecuteTemplate(c.Response(), "invoice_table", nil)
}

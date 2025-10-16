// Package controller
package controller

import (
	"fmt"
	"net/http"
	"strings"
	"vis_contas/internal/service"
	"vis_contas/utils"

	"github.com/labstack/echo/v4"
)

// LoadTable retorna os dados de acordo com todos os filtros
func LoadTable(c echo.Context) error {
	// Obter filtros
	categoria := strings.TrimSpace(c.QueryParam("categoria"))
	situacao := strings.TrimSpace(c.QueryParam("situacao"))
	year := strings.TrimSpace(c.QueryParam("year"))
	month := strings.TrimSpace(c.QueryParam("month"))
	empresa := strings.TrimSpace(c.QueryParam("empresa"))

	// Formata a Data a partir de Year e Month
	dataStr := utils.ConstruirData(year, month)

	faturasFiltradas, err := service.FilterFaturas(categoria, situacao, dataStr, empresa)
	if err != nil {
		fmt.Println("Erro ao obter Faturas Filtradas:", err)
		return err
	}

	// Preparar dados para o template
	dataMap := map[string]any{
		"Faturas": faturasFiltradas,
	}

	return c.Render(http.StatusOK, "invoice_table", dataMap)
}

func AlternarSituacao(c echo.Context) error {
	// Id da Fatura
	idStr := c.Param("id")

	// Filtros da Página
	categoria := strings.TrimSpace(c.FormValue("categoria"))
	situacaoFiltro := strings.TrimSpace(c.FormValue("situacao"))
	year := strings.TrimSpace(c.FormValue("year"))
	month := strings.TrimSpace(c.FormValue("month"))
	empresa := strings.TrimSpace(c.FormValue("empresa"))

	// Formata a Data a partir de Year e Month
	dataStr := utils.ConstruirData(year, month)

	// Atualizar dados da Fatura Selecionada
	err := service.AtualizarSituacao(idStr)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	faturasFiltradas, err := service.FilterFaturas(categoria, situacaoFiltro, dataStr, empresa)
	if err != nil {
		fmt.Println("Erro ao obter Faturas Filtradas:", err)
		return err
	}

	return c.Render(http.StatusOK, "invoice_table", map[string]any{"Faturas": faturasFiltradas})
}

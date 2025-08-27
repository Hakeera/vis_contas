// Package controller
package controller

import (
	"fmt"
	"net/http"
	"strings"
	"vis_contas/internal/service"

	"github.com/labstack/echo/v4"
)

// LoadTable retorna os dados de acordo com todos os filtros
func LoadTable(c echo.Context) error {
	// Ler CSV
	faturas, err := service.ReadFaturasCSV("data/csv/faturas.csv")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Erro ao ler arquivo csv")
	}

	// Obter filtros
	categoria := strings.TrimSpace(c.QueryParam("categoria"))
	situacao := strings.TrimSpace(c.QueryParam("situacao"))
	dataStr := strings.TrimSpace(c.QueryParam("data")) // TODO Ajustar Formatção de Data
	empresa := strings.TrimSpace(c.QueryParam("empresa")) 

	faturas, err = service.FilterFaturas(categoria, situacao, dataStr, empresa)
	if err != nil {
		fmt.Println("Erro ao obter Faturas Filtradas:", err)
		return err
	}
	// Preparar dados para o template
	dataMap := map[string]any{
		"Faturas": faturas,
	}

	fmt.Println("Dados:", empresa)
	return c.Render(http.StatusOK, "invoice_table", dataMap)
}

func DeixarPendente(c echo.Context) error {
	// Ler CSV
	faturas, err := service.ReadFaturasCSV("data/csv/faturas.csv")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Erro ao ler arquivo csv")
	}

	idStr := c.Param("id")
	situacao := "Pago"

	// Atualiza
	service.AtualizarSituacao(situacao, idStr)

	// Preparar dados para o template
	dataMap := map[string]any{
		"Faturas": faturas,
	}

	return c.Render(http.StatusOK, "invoice_table", dataMap)
}

func DeixarPago(c echo.Context) error {

	idStr := c.Param("id")
	situacao := "Pendente"

	// Atualiza
	faturas := service.AtualizarSituacao(situacao, idStr)

	// Preparar dados para o template
	dataMap := map[string]any{
		"Faturas": faturas,
	}

	return c.Render(http.StatusOK, "invoice_table", dataMap)
}

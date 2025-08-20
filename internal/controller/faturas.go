// Package controller
package controller

import (
	"fmt"
	"net/http"
	"strings"
	"vis_contas/service"

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

    faturas = service.FilterFaturas(faturas, categoria, situacao, dataStr)

    // Preparar dados para o template
    dataMap := map[string]any{
        "Faturas": faturas,
    }

    return c.Render(http.StatusOK, "invoice_table", dataMap)
}

func AlterarSituacao(c echo.Context) error {
	idStr := c.Param("id")

	fmt.Println(idStr) 
	// TODO SERVICE: Criar atualização de situação no banco de dados
	return nil 
}

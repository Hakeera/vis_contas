// Package controller
package controller

import (
	"log"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
)


func Home(c echo.Context) error {
	tmpl, err := template.ParseFiles("view/faturas.html",
					"view/template/base.html",
)

	if err != nil {
		log.Println("Erro ao carregar templates:", err)
		return c.String(http.StatusInternalServerError, "Erro ao carregar templates: "+err.Error())
	}

	// Executa o template base
	return tmpl.ExecuteTemplate(c.Response(), "faturas.html", nil)
	
}

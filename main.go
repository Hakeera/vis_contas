package main

import (
	"fmt"
	"io"
	"log"
	"text/template"
	"vis_contas/internal/routes"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer para Echo
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders tamplates with data
func (t *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	// Debug: verificar se o template existe
	tmpl := t.templates.Lookup(name)
	if tmpl == nil {
		log.Printf("❌ Template '%s' não encontrado!", name)
		// Listar todos os templates disponíveis
		for _, t := range t.templates.Templates() {
			log.Printf("📄 Template disponível: %s", t.Name())
		}
		return fmt.Errorf("template %s não encontrado", name)
	}

	log.Printf("✅ Renderizando template: %s", name)
	log.Printf("📊 Dados: %+v", data)

	return t.templates.ExecuteTemplate(w, name, data)
}

func main()  {
	// Configura o renderer de templates com as funções personalizadas
	renderer := &TemplateRenderer{
		templates: template.Must(
			template.New("").ParseGlob("view/**/*.html"),
		),
	}

	// Instância do Echo
	e := echo.New()
	e.Static("/static", "view/static")
	e.Renderer = renderer

	// Configurar rotas
	routes.SetUpRoutes(e)

	// Iniciar o servidor
	e.Logger.Fatal(e.Start(":1323"))
}

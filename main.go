// main.go
// Aplicação principal do sistema de controle de pagamentos.
// Responsável por inicializar o servidor Echo, carregar templates e conectar ao banco de dados.

package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"vis_contas/config"
	"vis_contas/internal/routes"
	"vis_contas/utils"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer implementa o renderer de templates HTML para o Echo.
type TemplateRenderer struct {
	templates *template.Template
}

// Render executa o template solicitado e escreve a saída no Writer.
func (t *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	tmpl := t.templates.Lookup(name)
	if tmpl == nil {
		log.Printf("❌ Template '%s' não encontrado!", name)
		for _, t := range t.templates.Templates() {
			log.Printf("📄 Template disponível: %s", t.Name())
		}
		return fmt.Errorf("template %s não encontrado", name)
	}

	log.Printf("✅ Renderizando template: %s", name)
	log.Printf("📊 Dados: %+v", data)

	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("❌ ERRO na execução do template '%s': %v", name, err)
		return fmt.Errorf("erro ao executar template %s: %v", name, err)
	}

	log.Printf("✅ Template '%s' executado com sucesso!", name)
	return nil
}

func main() {
	// Carrega variáveis de ambiente e inicializa o banco
	config.LoadEnv()
	config.InitDB()

	// Teste de conexão com o banco
	db := config.GetDB()
	if db != nil {
		log.Println("✅ Banco de dados conectado com sucesso!")
		if sqlDB, err := db.DB(); err == nil {
			if err := sqlDB.Ping(); err == nil {
				log.Println("✅ Ping no banco OK!")
			} else {
				log.Printf("❌ Erro no ping: %v", err)
			}
		}
	} else {
		log.Println("❌ Banco de dados é nil!")
	}

	tmpl := template.New("").Funcs(utils.TemplateFunctions)
	tmpl = template.Must(tmpl.ParseGlob("view/**/*.html"))

	renderer := &TemplateRenderer{
		templates: tmpl,
	}

	// Inicialização do servidor Echo
	e := echo.New()
	e.Static("/static", "view/static")
	e.Renderer = renderer

	// Configuração das rotas da aplicação
	routes.SetUpRoutes(e)

	// Inicia o servidor na porta 8080
	log.Println("🚀 Servidor iniciando na porta :1323")
	e.Logger.Fatal(e.Start(":8080"))
}

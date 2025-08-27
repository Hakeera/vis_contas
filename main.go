package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"vis_contas/config"
	"vis_contas/internal/routes"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer para Echo
type TemplateRenderer struct {
	templates *template.Template
}

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
    
    // 🎯 Capturar o erro da execução
    err := t.templates.ExecuteTemplate(w, name, data)
    if err != nil {
        log.Printf("❌ ERRO na execução do template '%s': %v", name, err)
        return fmt.Errorf("erro ao executar template %s: %v", name, err)
    }
    
    log.Printf("✅ Template '%s' executado com sucesso!", name)
    return nil
}

func main() {

	// Carregar Variáveis de Ambiente
	config.LoadEnv()
	
	// Inicializa o banco de dados
	config.InitDB()

	// Verifica se o banco está funcionando
	db := config.GetDB()
	if db != nil {
		log.Println("✅ Banco de dados conectado com sucesso!")
		// Teste simples de conexão
		sqlDB, err := db.DB()
		if err == nil {
			if err := sqlDB.Ping(); err == nil {
				log.Println("✅ Ping no banco OK!")
			} else {
				log.Printf("❌ Erro no ping: %v", err)
			}
		}
	} else {
		log.Println("❌ Banco de dados é nil!")
	}

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
	log.Println("🚀 Servidor iniciando na porta :1323")
	e.Logger.Fatal(e.Start(":1323"))
}

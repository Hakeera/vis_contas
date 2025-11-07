// internal/routes/routes.go
// Define as rotas e middlewares principais da aplicação Echo.

package routes

import (
	"vis_contas/internal/auth"
	"vis_contas/internal/controller"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetUpRoutes configura as rotas HTTP e middlewares globais da aplicação.
//
// Rotas públicas:
//   - GET  /user   → Exibe página de login
//   - POST /login  → Autentica usuário e gera token JWT
//
// Rotas protegidas (JWT):
//   - GET  /           → Página inicial
//   - GET  /load_table → Atualiza tabela principal via HTMX
//   - PUT  /alternar_sit/:id → Alterna situação de uma fatura
//   - POST /logout     → Finaliza sessão do usuário
func SetUpRoutes(e *echo.Echo) {

	// Middlewares globais
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// 👥 Rotas públicas
	e.GET("/user", controller.LoginPage)
	e.POST("/login", controller.Login)

	// Grupo protegido por autenticação JWT
	r := e.Group("")
	r.Use(auth.RequireAuth)

	r.GET("/", controller.Home)
	r.GET("/load_table", controller.LoadTable)
	r.PUT("/alternar_sit/:id", controller.AlternarSituacao)
	r.POST("/logout", controller.Logout)
}

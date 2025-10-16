package controller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"vis_contas/internal/auth"
	"vis_contas/internal/service"

	"github.com/labstack/echo/v4"
)

// TODO Adaptar dinâmica de templates com c.Render
// TODO Configurar cookies e token de autenticação

// UserRequest Estrutura para receber os dados do formulário de login/registro
type UserRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

// LoginPage GET / - Serve a página de login
func LoginPage(c echo.Context) error {
	tmpl, err := template.ParseFiles("view/login.html")
	if err != nil {
		log.Println("Erro ao carregar templates:", err)
		return c.String(http.StatusInternalServerError, "Erro ao carregar templates: "+err.Error())
	}

	// Executa o template base
	return tmpl.ExecuteTemplate(c.Response(), "login.html", nil)
}

// Register POST /register - Registra um novo usuário
func Register(c echo.Context) error {
	var req UserRequest
	if err := c.Bind(&req); err != nil {
		return c.HTML(http.StatusBadRequest, `
			<div class="alert alert-error">
				Dados inválidos. Verifique os campos e tente novamente.
			</div>
		`)
	}

	// Validações básicas
	if strings.TrimSpace(req.Username) == "" {
		return c.HTML(http.StatusBadRequest, `
			<div class="alert alert-error">
				Nome de usuário é obrigatório.
			</div>
		`)
	}

	if len(req.Password) < 6 {
		return c.HTML(http.StatusBadRequest, `
			<div class="alert alert-error">
				A senha deve ter pelo menos 6 caracteres.
			</div>
		`)
	}

	user, err := service.CreateUser(req.Username, req.Password)
	if err != nil {
		return c.HTML(http.StatusBadRequest, `
			<div class="alert alert-error">
				`+err.Error()+`
			</div>
		`)
	}

	return c.HTML(http.StatusCreated, `
		<div class="alert alert-success">
			Usuário "`+template.HTMLEscapeString(user.Username)+`" criado com sucesso! 
			<br>Agora você pode fazer login.
		</div>
	`)
}

// Login POST /login - Autentica um usuário
func Login(c echo.Context) error {
	var req UserRequest
	if err := c.Bind(&req); err != nil {
		fmt.Println("DADOs INVÁLIDOs")
		return c.String(http.StatusBadRequest, "Dados inválidos")
	}

	user, err := service.AutenticarUsuario(req.Username, req.Password)
	if err != nil {
		// Preparar dados para o template
		dataMap := map[string]any{
			"Mensagem": "Usuário ou Senha incorretos!",
		}
		return c.Render(http.StatusUnauthorized, "login.html", dataMap)
	}

	// Gera token JWT
	token, err := auth.GerarToken(user.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Erro ao gerar token")
	}

	// Salva token em cookie seguro
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   600, // 1h
	})

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
}

// Logout POST /logout - Faz logout do usuário
func Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	return c.Redirect(http.StatusSeeOther, "/user")
}

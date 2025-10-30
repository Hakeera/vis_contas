package controller

import (
	"fmt"
	"net/http"

	"vis_contas/internal/auth"
	"vis_contas/internal/service"

	"github.com/labstack/echo/v4"
)

// TODO: Adaptar dinâmica de templates com c.Render

// UserRequest Estrutura para receber os dados do formulário de login/registro
type UserRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

// LoginPage GET / - Serve a página de login
func LoginPage(c echo.Context) error {

	// Executa o template base
	return c.Render(http.StatusOK, "login", nil)
}

// Login POST /login - Autentica um usuário
func Login(c echo.Context) error {
	var req UserRequest
	if err := c.Bind(&req); err != nil {
		fmt.Println("DADOs INVÁLIDOS")
		return c.String(http.StatusBadRequest, "Dados inválidos")
	}

	user, err := service.AutenticarUsuario(req.Username, req.Password)
	if err != nil {
		return c.String(http.StatusOK, "Usuário ou Senha incorretos!")
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
		MaxAge:   3600, // 1h
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
	c.Response().Header().Set("HX-Redirect", "/user")
	return c.NoContent(http.StatusOK)
}

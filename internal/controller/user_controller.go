package controller

import (
	"net/http"

	"vis_contas/internal/service"

	"github.com/labstack/echo/v4"
)

// UserRequest Estrutura para receber os dados do formulário de login/registro
type UserRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

// Register POST /register
func Register(c echo.Context) error {
	var req UserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "dados inválidos"})
	}

	user, err := service.CreateUser(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"message": "usuário criado com sucesso",
		"user":    user.Username,
	})
}

// Login POST /login
func Login(c echo.Context) error {
	var req UserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "dados inválidos"})
	}

	user, err := service.AutenticarUsuario(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// Aqui futuramente podemos adicionar JWT ou Cookie de sessão
	return c.JSON(http.StatusOK, map[string]any{
		"message": "login realizado com sucesso",
		"user":    user.Username,
	})
}

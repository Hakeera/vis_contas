// internal/auth/jwt.go
// Implementa autenticação baseada em JWT para a aplicação.

package auth

import (
	"net/http"
	"time"
	"vis_contas/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GerarToken cria um token JWT com ID do usuário e expiração de 2 horas.
// Retorna o token assinado ou erro em caso de falha.
func GerarToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWTSecret)
}

// RequireAuth é um middleware que valida o token JWT armazenado no cookie "token".
// Redireciona para /user se o token for inválido ou ausente.
// Define "user_id" no contexto quando o token é válido.
func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/user")
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.ErrUnauthorized
			}
			return config.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			return c.Redirect(http.StatusSeeOther, "/user")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
		}

		return next(c)
	}
}

// config/config.go
// Pacote responsável pelo carregamento de variáveis de ambiente da aplicação.

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// JWTSecret contém a chave secreta utilizada na geração e validação de tokens JWT.
var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

// LoadEnv carrega variáveis do arquivo .env, caso exista.
// Se o arquivo não for encontrado, utiliza variáveis de ambiente do sistema.
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Aviso: .env não encontrado, usando variáveis de ambiente do sistema.")
	}
}

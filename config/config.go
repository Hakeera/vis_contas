package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

// LoadEnv carrega variáveis do arquivo .env
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Não foi possível carregar o arquivo .env, usando variáveis de ambiente existentes")
	}
}

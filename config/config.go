package config 

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv carrega variáveis do arquivo .env
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Não foi possível carregar o arquivo .env, usando variáveis de ambiente existentes")
	}
}

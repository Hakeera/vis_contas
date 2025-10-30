// Package config
package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"vis_contas/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// DB é a instância global de conexão com o banco de dados
	DB *gorm.DB
	// once garante que a inicialização do banco ocorra apenas uma vez
	once sync.Once
)

// InitDB inicializa a conexão com o banco de dados PostgreSQL
// utilizando as variáveis de ambiente para configuração
// A função também realiza a migração automática das tabelas necessárias
// Esta função é thread-safe e garante que a inicialização ocorra apenas uma vez
func InitDB() {
	once.Do(func() {
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbSchema := os.Getenv("DB_SCHEMA")
		sslmode := os.Getenv("DB_SSLMODE")

		// Define o search_path para o schema financeiro
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s search_path=%s",
			dbHost, dbUser, dbPassword, dbName, dbPort, sslmode, dbSchema,
		)

		var err error
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("❌ Erro ao conectar ao banco de dados: %v", err)
		}

		log.Println("✅ Conectado ao banco de dados com sucesso!")

		// Migrar as tabelas (opcional se já existir via init.sql)
		if err := DB.AutoMigrate(
			&model.Fatura{},
			&model.User{},
		); err != nil {
			log.Fatalf("❌ Erro ao migrar banco: %v", err)
		}
	})
}

// GetDB retorna a instância do banco de dados inicializada
// Causa erro fatal se o banco de dados não foi inicializado previamente
// @return *gorm.DB Instância de conexão com o banco de dados
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Banco de dados não foi inicializado!")
	}
	return DB
}

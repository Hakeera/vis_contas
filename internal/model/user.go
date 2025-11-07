// internal/model/user.go
// Define a entidade User e seu mapeamento para a tabela financeiro.users.

package model

import "time"

// User representa um usuário do sistema.
// Armazena credenciais, papel (role) e timestamps de criação/atualização.
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"unique;not null" json:"username"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         string    `gorm:"default:user" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName define o nome completo da tabela no banco de dados.
func (User) TableName() string {
	return "financeiro.users"
}

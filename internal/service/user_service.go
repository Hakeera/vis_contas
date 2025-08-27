package service

import (
	"errors"
	"vis_contas/config"
	"vis_contas/internal/model"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

// CreateUser registra um novo usuário no banco
func CreateUser(username, password string) (*model.User, error) {

	db := config.GetDB()

	// Verifica se já existe usuário com esse username
	var existing model.User
	if err := db.Where("username = ?", username).First(&existing).Error; err == nil {
		return nil, errors.New("usuário já existe")
	}

	// Gera o hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Cria usuário
	user := &model.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Role:         "user",
	}

	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// AutenticarUsuario faz login verificando username e senha
func AutenticarUsuario(username, password string) (*model.User, error) {
	db := config.GetDB()

	var user model.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}

	// Verifica a senha
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("senha incorreta")
	}

	return &user, nil
}

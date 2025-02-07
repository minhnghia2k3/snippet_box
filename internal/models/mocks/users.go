package mocks

import (
	"github.com/minhnghia2k3/snippet_box/internal/models"
	"time"
)

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@gmail.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "real@gmail.com" && password == "pa$$word" {
		return 1, nil
	}
	return 0, models.ErrNoRecord
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) Get(id int) (*models.User, error) {
	user := models.User{
		ID:           id,
		Name:         "test",
		Email:        "test@gmail.com",
		HashPassword: []byte("pa$$word"),
		Created:      time.Now(),
	}
	return &user, nil
}

func (m *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	if id == 1 {
		if currentPassword != "pa$$word" {
			return models.ErrInvalidCredentials
		}
		return nil
	}
	return models.ErrNoRecord
}

package mock

import (
	"time"

	"github.com/katarzynakawala/diffs-app/pkg/models"
)

var mockUser = &models.User{
	ID: 1,
	Name: "Kasia",
	Email: "kasia@kasia.com",
	Created: time.Now(),
	Active: true,
}

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "example@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil	
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	switch email {
	case "kasia@kasia.com":
		return 1, nil
	default:
		return 0, models.ErrInvalidCredentials	
	}
}

func (m *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord	
	}
}
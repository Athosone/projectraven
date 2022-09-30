package domain

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID
	Email       string
	IdpProvider string
}

type UserRepository interface {
	FindById(id uuid.UUID) (*User, error)
	FindByIdpProviderAndIdpId(idpProvider string, idpId string) (*User, error)
	FindByEmail(email string) (*User, error)
	Save(user *User) error
}

package app

import (
	"errors"

	"github.com/athosone/projectraven/tracking/internal/domain"
)

type RegisterUserCommandHandler struct {
	repo domain.UserRepository
}

func NewRegisterUserCommandHandler(repo domain.UserRepository) (*RegisterUserCommandHandler, error) {
	return &RegisterUserCommandHandler{repo: repo}, nil
}

func (h *RegisterUserCommandHandler) Handle(command RegisterUserCommand) error {
	_, err := h.repo.FindByIdpProviderAndIdpId(command.IdpProvider, command.IdpId)
	if err == nil {
		return ErrUserAlreadyRegistered
	}
	user := &domain.User{
		Email:       command.Email,
		IdpProvider: command.IdpProvider,
	}
	return h.repo.Save(user)
}

type RegisterUserCommand struct {
	Email       string
	IdpProvider string
	IdpId       string
}

var ErrUserAlreadyRegistered = errors.New("user already registered")

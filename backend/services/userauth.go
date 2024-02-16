package services

import (
	"context"

	// bean "github.com/retail-ai-inc/bean/v2/v2/trace"

	"backend/models"
	"backend/packages/global"
	"backend/repositories"
)

type UserauthService interface {
	SignUp(ctx context.Context, userToAdd models.User) error
	SignIn(ctx context.Context, userEmail string, userPassword string) (*models.User, error)
}

type userauthService struct {
	userauthRepository repositories.UserauthRepository
}

func NewUserauthService(userauthRepo repositories.UserauthRepository) *userauthService {
	return &userauthService{
		userauthRepository: userauthRepo,
	}
}

func (service *userauthService) SignUp(ctx context.Context, userToAdd models.User) error {
	err := service.userauthRepository.UserSignUp(ctx, userToAdd)
	return err
}

func (service *userauthService) SignIn(ctx context.Context, userEmail string, userPassword string) (*models.User, error) {
	user, err := service.userauthRepository.AuthenticateUser(ctx, userEmail, userPassword)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, global.ErrInvalidUserCredentials
	}

	return user, err
}

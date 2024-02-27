package services

import (
	"backend/models"
	"backend/packages/global"
	"backend/repositories"
	"context"
)

type AdminAuthService interface {
	SignIn(ctx context.Context, userEmail string, userPassword string) (*models.Admin, error)
}

type adminAuthService struct {
	adminAuthRepository repositories.AdminAuthRepository
}

func NewAdminAuthService(adminAuthRepo repositories.AdminAuthRepository) *adminAuthService {
	return &adminAuthService{
		adminAuthRepository: adminAuthRepo,
	}
}

func (service *adminAuthService) SignIn(ctx context.Context, userEmail string, userPassword string) (*models.Admin, error) {
	admin, err := service.adminAuthRepository.AuthenticateAdmin(ctx, userEmail, userPassword)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, global.ErrInvalidUserCredentials
	}

	return admin, err
}

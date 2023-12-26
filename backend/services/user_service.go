package services

import (
	"context"
	"time"

	// "github.com/retail-ai-inc/bean/trace"
	"backend/repositories"
)

type UserService interface {
	IsUserApiQuotaRemaining(ctx context.Context, userIp string) (bool, error)
	GetUserTTL(ctx context.Context, uerIp string) (time.Duration, error)
	DecrementApiQouta(ctx context.Context, userIp string) error
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *userService {
	return &userService{
		userRepository: userRepo,
	}
}

func (service *userService) IsUserApiQuotaRemaining(ctx context.Context, userIp string) (bool, error) {
	quota, err := service.userRepository.GetUserQuota(ctx, userIp)
	if err != nil {
		panic(err)
	}

	if quota > 0 {
		return true, nil
	}

	return false, nil
}

func (service *userService) GetUserTTL(ctx context.Context, userIp string) (time.Duration, error) {
	ttl, err := service.userRepository.GetUserTTL(ctx, userIp)
	if err != nil {
		panic(err)
	}

	return ttl, nil
}

func (service *userService) DecrementApiQouta(ctx context.Context, userIp string) error {
	service.userRepository.DecrementQuota(ctx, userIp)
	return nil
}

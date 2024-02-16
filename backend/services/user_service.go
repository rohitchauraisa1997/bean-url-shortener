package services

import (
	"context"
	"time"

	// bean "github.com/retail-ai-inc/bean/v2/trace"
	"backend/repositories"
)

type UserService interface {
	IsUserApiQuotaRemaining(ctx context.Context, userIp string) (bool, error)
	GetUserTTL(ctx context.Context, uerIp string) (time.Duration, error)
	AddUrlToUserKeyAndDecrementApiQouta(ctx context.Context, userIp string, shortenedUrl string) error
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

func (service *userService) AddUrlToUserKeyAndDecrementApiQouta(ctx context.Context, userIp string, shortenedUrl string) error {
	service.userRepository.AddUrlToUserKeyAndDecrementApiQouta(ctx, userIp, shortenedUrl)
	return nil
}

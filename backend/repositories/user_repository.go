package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/retail-ai-inc/bean/v2/store/redis"
)

const API_QUOTA = 100

type UserRepository interface {
	GetUserQuota(c context.Context, userId string) (uint64, error)
	GetUserTTL(c context.Context, userId string) (time.Duration, error)
	AddUrlToUserKeyAndDecrementApiQouta(c context.Context, userId string, shortenedUrl string) error
}

type userRepository struct {
	userRepo redis.MasterCache
}

func NewUserRepository(userRepo redis.MasterCache) *userRepository {
	return &userRepository{userRepo}
}

func (r *userRepository) GetUserQuota(c context.Context, userId string) (uint64, error) {
	// first check if quoata for user exists or not.
	// if key for the user doesnt exist it means that the quota was not assigned to the user, assign a quota to that user.
	// if key for the user exists then return the quota for that user.
	userKeyToSearchFor := fmt.Sprintf("User%s", userId)
	keyExists, err := r.userRepo.KeyExists(c, userKeyToSearchFor)
	if err != nil {
		panic(err)
	}

	if keyExists {
		currQuota, err := r.userRepo.HGet(c, userKeyToSearchFor, "currQuota")
		if err != nil {
			panic(err)
		}
		currQuotaUint, err := strconv.Atoi(currQuota)
		if err != nil {
			panic(err)
		}
		return uint64(currQuotaUint), nil
	} else {
		// set the user quota for 24 hours.
		err = r.userRepo.HSet(c, userKeyToSearchFor, "currQuota", API_QUOTA)
		if err != nil {
			panic(err)
		}
		err = r.userRepo.Expire(c, userKeyToSearchFor, time.Hour*24)
		if err != nil {
			panic(err)
		}
		return API_QUOTA, nil
	}

}

func (r *userRepository) GetUserTTL(c context.Context, userId string) (time.Duration, error) {
	userKeyToSearchFor := fmt.Sprintf("User%s", userId)
	keyExists, err := r.userRepo.KeyExists(c, userKeyToSearchFor)
	if err != nil {
		panic(err)
	}
	if keyExists {
		ttl, err := r.userRepo.Ttl(c, userKeyToSearchFor)
		if err != nil {
			panic(err)
		}
		return ttl, nil
	}
	return 0, nil
}

func (r *userRepository) AddUrlToUserKeyAndDecrementApiQouta(c context.Context, userId string, shortenedUrl string) error {
	userKeyToSearchFor := fmt.Sprintf("User%s", userId)
	currQuota, err := r.userRepo.HGet(c, userKeyToSearchFor, "currQuota")
	if err != nil {
		panic(err)
	}

	updatedCurrQuota, err := strconv.Atoi(currQuota)
	if err != nil {
		panic(err)
	}

	updatedCurrQuota -= 1
	err = r.userRepo.HSet(c, userKeyToSearchFor, "currQuota", updatedCurrQuota)
	if err != nil {
		panic(err)
	}

	shortenedUrlsCreatedByUser, err := r.userRepo.HGet(c, userKeyToSearchFor, "urlsList")
	if err != nil {
		panic(err)
	}

	if shortenedUrlsCreatedByUser == "" {
		shortenedUrlsCreatedByUserList := []string{shortenedUrl}
		marshaledshortenedUrlsCreatedByUserList, _ := json.Marshal(shortenedUrlsCreatedByUserList)
		err = r.userRepo.HSet(c, userKeyToSearchFor, "urlsList", marshaledshortenedUrlsCreatedByUserList)
		if err != nil {
			panic(err)
		}
	} else {
		var shortenedUrlsCreatedByUserList []string
		json.Unmarshal([]byte(shortenedUrlsCreatedByUser), &shortenedUrlsCreatedByUserList)
		shortenedUrlsCreatedByUserList = append(shortenedUrlsCreatedByUserList, shortenedUrl)
		marshaledshortenedUrlsCreatedByUserList, _ := json.Marshal(shortenedUrlsCreatedByUserList)
		err = r.userRepo.HSet(c, userKeyToSearchFor, "urlsList", marshaledshortenedUrlsCreatedByUserList)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

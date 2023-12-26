package repositories

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/retail-ai-inc/bean/dbdrivers"
)

const API_QUOTA = 100

type UserRepository interface {
	GetUserQuota(c context.Context, userIp string) (uint64, error)
	GetUserTTL(c context.Context, userIp string) (time.Duration, error)
	DecrementQuota(c context.Context, userIp string) error
}

type userRepository struct {
	client *dbdrivers.RedisDBConn
}

func NewUserRepository(client *dbdrivers.RedisDBConn) *userRepository {
	return &userRepository{client}
}

func (r *userRepository) GetUserQuota(c context.Context, userIp string) (uint64, error) {
	// first check if quoata for user exists or not.
	// if key for the user doesnt exist it means that the quota was not assigned to the user, assign a quota to that user.
	// if key for the user exists then return the quota for that user.
	// userIp = "localhost"
	userKeyToSearchFor := fmt.Sprintf("User%s", userIp)
	keyExists, err := dbdrivers.RedisIsKeyExists(c, r.client, userKeyToSearchFor)
	if err != nil {
		panic(err)
	}

	if keyExists {
		fmt.Println("User exists in db")
		currQuota, err := dbdrivers.RedisHGet(c, r.client, userKeyToSearchFor, "currQuota")
		if err != nil {
			panic(err)
		}
		currQuotaUint, err := strconv.Atoi(currQuota)
		if err != nil {
			panic(err)
		}
		return uint64(currQuotaUint), nil
	} else {
		fmt.Println("User NOTNOTNOTNOT exists in db")
		// set the user quota for 24 hours.
		err = dbdrivers.RedisHSet(c, r.client, userKeyToSearchFor, "currQuota", API_QUOTA, time.Hour*24)
		if err != nil {
			panic(err)
		}
		return API_QUOTA, nil
	}

}

func (r *userRepository) GetUserTTL(c context.Context, userIp string) (time.Duration, error) {
	userKeyToSearchFor := fmt.Sprintf("User%s", userIp)
	keyExists, err := dbdrivers.RedisIsKeyExists(c, r.client, userKeyToSearchFor)
	if err != nil {
		panic(err)
	}
	if keyExists {
		ttl, err := dbdrivers.RedisGetTTL(c, r.client, userKeyToSearchFor)
		if err != nil {
			panic(err)
		}
		return ttl, nil
	}
	return 0, nil
}

func (r *userRepository) DecrementQuota(c context.Context, userIp string) error {
	userKeyToSearchFor := fmt.Sprintf("User%s", userIp)
	currQuota, err := dbdrivers.RedisHGet(c, r.client, userKeyToSearchFor, "currQuota")
	if err != nil {
		panic(err)
	}
	updatedCurrQuota, err := strconv.Atoi(currQuota)
	if err != nil {
		panic(err)
	}
	updatedCurrQuota -= 1
	err = dbdrivers.RedisHSet(c, r.client, userKeyToSearchFor, "currQuota", updatedCurrQuota, 0)
	if err != nil {
		panic(err)
	}
	return nil
}

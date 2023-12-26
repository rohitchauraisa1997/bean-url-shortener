package repositories

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/retail-ai-inc/bean/dbdrivers"
)

type UrlRepository interface {
	GetNewShortenedUrl(c context.Context) (string, error)
	CreateShortenedUrl(c context.Context, shortenedUrl string, url string, userIp string, ttlInMins time.Duration) (string, error)
	GetUrlViaShortenedUrl(c context.Context, shortenedUrl string) (string, error)
	IncreaseVisitCounter(c context.Context, shortenedUrl string) (uint64, error)
	GetAllShortenedUrlKeys(c context.Context) ([]string, error)
	GetUrlForShortenedKeys(c context.Context, shortenedKeys []string) (map[string]string, error)
	GetTtlForShortenedKeys(c context.Context, shortenedKeys []string) (map[string]uint64, error)
	GetUrlHitsForShortenedKeys(c context.Context, shortenedKeys []string) (map[string]uint64, error)
	GetCreatedByForShortenedKeys(c context.Context, shortenedKeys []string) (map[string]string, error)
}

type urlRepository struct {
	client *dbdrivers.RedisDBConn
}

func NewUrlRepository(client *dbdrivers.RedisDBConn) *urlRepository {
	return &urlRepository{client}
}

func (r *urlRepository) GetNewShortenedUrl(c context.Context) (string, error) {
	fmt.Println("GetNewShortenedUrl triggered!!")
	var shortenedUrl string

	allKeys, err := dbdrivers.RedisGetKeys(c, r.client, "*")
	if err != nil {
		panic(err)
	}

	for {
		shortenedUrl = uuid.New().String()[:6]
		// Check if the generated shortenedUrl already exists in allKeys
		found := false
		for _, redisKey := range allKeys {
			if shortenedUrl == redisKey {
				found = true
				break
			}
		}

		if !found {
			break // Found a unique shortenedUrl, break the loop
		}
	}
	// dbdrivers.RedisHSet(c, r.client, shortenedUrl)

	return shortenedUrl, nil
}

func (r *urlRepository) CreateShortenedUrl(c context.Context, shortenedUrl string, urlToShorten string, userIp string, ttlInMins time.Duration) (string, error) {
	fmt.Println("CreateShortenedUrl triggered!!")
	err := dbdrivers.RedisHSet(c, r.client, shortenedUrl, "url", urlToShorten, 0)
	if err != nil {
		panic(err)
	}
	createdAt := time.Now().Unix()
	// adding createdAt and toBeDeletedAt at time of creation for preventing ttl calls for each key.
	toBeDeletedAt := time.Now().Unix() + int64((ttlInMins/time.Nanosecond/time.Minute)*60)
	err = dbdrivers.RedisHSet(c, r.client, shortenedUrl, "createdAt", createdAt, 0)
	if err != nil {
		panic(err)
	}
	err = dbdrivers.RedisHSet(c, r.client, shortenedUrl, "toBeDeletedAt", toBeDeletedAt, 0)
	if err != nil {
		panic(err)
	}
	err = dbdrivers.RedisHSet(c, r.client, shortenedUrl, "createdBy", userIp, 0)
	if err != nil {
		panic(err)
	}
	err = dbdrivers.RedisHSet(c, r.client, shortenedUrl, "urlHits", 0, ttlInMins)
	if err != nil {
		panic(err)
	}
	return shortenedUrl, nil
}

func (r *urlRepository) GetUrlViaShortenedUrl(c context.Context, shortenedUrl string) (string, error) {
	actualUrl, err := dbdrivers.RedisHGet(c, r.client, shortenedUrl, "url")
	if err != nil {
		panic(err)
	}
	actualUrl = strings.Trim(actualUrl, "\"")
	return actualUrl, nil
}

func (r *urlRepository) IncreaseVisitCounter(c context.Context, shortenedUrl string) (uint64, error) {
	currUrlHits, err := dbdrivers.RedisHGet(c, r.client, shortenedUrl, "urlHits")
	if err != nil {
		panic(err)
	}
	currUrlHitsUint, err := strconv.Atoi(currUrlHits)
	if err != nil {
		panic(err)
	}
	err = dbdrivers.RedisHSet(c, r.client, shortenedUrl, "urlHits", currUrlHitsUint+1, 0)
	if err != nil {
		panic(err)
	}
	return uint64(currUrlHitsUint) + 1, nil
}

func (r *urlRepository) GetAllShortenedUrlKeys(c context.Context) ([]string, error) {
	allKeys, err := dbdrivers.RedisGetKeys(c, r.client, "*")
	if err != nil {
		panic(err)
	}
	onlyShortenedUrlKeys := []string{}
	for _, key := range allKeys {
		if !strings.Contains(key, "User") {
			onlyShortenedUrlKeys = append(onlyShortenedUrlKeys, key)
		}
	}
	return onlyShortenedUrlKeys, nil
}

func (r *urlRepository) GetUrlForShortenedKeys(c context.Context, shortenedKeys []string) (map[string]string, error) {
	var redisKeysWithField = make(map[string]string)
	for i := 0; i < len(shortenedKeys); i++ {
		redisKeysWithField[shortenedKeys[i]] = "url"
	}

	return dbdrivers.RedisHgets(c, r.client, redisKeysWithField)
}

func (r *urlRepository) GetTtlForShortenedKeys(c context.Context, shortenedKeys []string) (map[string]uint64, error) {
	var redisKeysWithField = make(map[string]string)
	for i := 0; i < len(shortenedKeys); i++ {
		redisKeysWithField[shortenedKeys[i]] = "toBeDeletedAt"
	}

	shoretenedKeyWithDeletedAt, err := dbdrivers.RedisHgets(c, r.client, redisKeysWithField)
	if err != nil {
		panic(err)
	}

	var shortenedKeysWithTtl = make(map[string]uint64)
	for skey, deletedAtStr := range shoretenedKeyWithDeletedAt {
		toBeDeletedAt, err := strconv.ParseUint(deletedAtStr, 10, 64)
		if err != nil {
			panic(err)
		}
		shortenedKeysWithTtl[skey] = toBeDeletedAt
	}

	return shortenedKeysWithTtl, nil
}

func (r *urlRepository) GetUrlHitsForShortenedKeys(c context.Context, shortenedKeys []string) (map[string]uint64, error) {
	var redisKeysWithField = make(map[string]string)
	for i := 0; i < len(shortenedKeys); i++ {
		redisKeysWithField[shortenedKeys[i]] = "urlHits"
	}

	shoretenedKeyWithUrlHits, err := dbdrivers.RedisHgets(c, r.client, redisKeysWithField)
	if err != nil {
		panic(err)
	}

	var shortenedKeysWithUrlHits = make(map[string]uint64)
	for skey, urlHitsStr := range shoretenedKeyWithUrlHits {
		urlHits, err := strconv.ParseUint(urlHitsStr, 10, 64)
		if err != nil {
			panic(err)
		}
		shortenedKeysWithUrlHits[skey] = urlHits
	}

	return shortenedKeysWithUrlHits, nil
}

func (r *urlRepository) GetCreatedByForShortenedKeys(c context.Context, shortenedKeys []string) (map[string]string, error) {
	var redisKeysWithField = make(map[string]string)
	for i := 0; i < len(shortenedKeys); i++ {
		redisKeysWithField[shortenedKeys[i]] = "createdBy"
	}

	return dbdrivers.RedisHgets(c, r.client, redisKeysWithField)
}

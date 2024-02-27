package repositories

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/retail-ai-inc/bean/v2/store/redis"
	"github.com/spf13/viper"
)

type UrlRepository interface {
	GetNewShortenedUrl(c context.Context) (string, error)
	CreateShortenedUrl(c context.Context, shortenedUrl string, url string, userIp string, ttlInMins time.Duration) (string, error)
	GetUrlViaShortenedUrl(c context.Context, shortenedUrl string) (string, error)
	IncreaseVisitCounter(c context.Context, shortenedUrl string) (uint64, error)
	GetAllShortenedUrlKeysForAllUsers(c context.Context) ([]string, error)
	GetAllShortenedUrlKeysForUser(c context.Context, userId string) ([]string, error)
	GetUrlForShortenedKeys(c context.Context, shortenedKeys []string) (map[string]string, error)
	GetTtlForShortenedKeys(c context.Context, shortenedKeys []string) (map[string]uint64, error)
	GetUrlHitsForShortenedKeys(c context.Context, shortenedKeys []string) (map[string]uint64, error)
	GetCreatedByForShortenedKeys(c context.Context, shortenedKeys []string) (map[string]string, error)
}

type urlRepository struct {
	urlRepo redis.MasterCache
}

func NewUrlRepository(urlRepo redis.MasterCache) *urlRepository {
	return &urlRepository{urlRepo}
}

func (r *urlRepository) GetNewShortenedUrl(c context.Context) (string, error) {
	var shortenedUrl string

	allKeys, err := r.urlRepo.Keys(c, "*")
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

	return shortenedUrl, nil
}

func (r *urlRepository) CreateShortenedUrl(c context.Context, shortenedUrl string, urlToShorten string, userId string, ttlInMins time.Duration) (string, error) {
	createdAt := time.Now().Unix()
	toBeDeletedAt := time.Now().Unix() + int64((ttlInMins/time.Nanosecond/time.Minute)*60)

	urlMap := map[string]interface{}{
		"url":           urlToShorten,
		"createdAt":     createdAt,
		"toBeDeletedAt": toBeDeletedAt,
		"createdBy":     userId,
		"urlHits":       0,
	}
	err := r.urlRepo.HSet(c, shortenedUrl, urlMap)
	if err != nil {
		panic(err)
	}

	err = r.urlRepo.Expire(c, shortenedUrl, ttlInMins)
	if err != nil {
		panic(err)
	}

	return shortenedUrl, nil
}

func (r *urlRepository) GetUrlViaShortenedUrl(c context.Context, shortenedUrl string) (string, error) {
	actualUrl, err := r.urlRepo.HGet(c, shortenedUrl, "url")
	if err != nil {
		panic(err)
	}
	actualUrl = strings.Trim(actualUrl, "\"")
	return actualUrl, nil
}

func (r *urlRepository) IncreaseVisitCounter(c context.Context, shortenedUrl string) (uint64, error) {
	currUrlHits, err := r.urlRepo.HGet(c, shortenedUrl, "urlHits")
	if err != nil {
		panic(err)
	}
	currUrlHitsUint, err := strconv.Atoi(currUrlHits)
	if err != nil {
		panic(err)
	}
	err = r.urlRepo.HSet(c, shortenedUrl, "urlHits", currUrlHitsUint+1)
	if err != nil {
		panic(err)
	}
	return uint64(currUrlHitsUint) + 1, nil
}

func (r *urlRepository) GetAllShortenedUrlKeysForAllUsers(c context.Context) ([]string, error) {
	allUserKeys, err := r.urlRepo.Keys(c, "User*")
	if err != nil {
		panic(err)
	}

	var redisKeysWithField = make(map[string]string)
	for i := 0; i < len(allUserKeys); i++ {
		allUserKeys[i] = strings.TrimPrefix(allUserKeys[i], viper.GetString("database.redis.prefix")+"_")
		redisKeysWithField[allUserKeys[i]] = "urlsList"
	}

	shortenedUrlsCreatedByUsers, err := r.urlRepo.HGets(c, redisKeysWithField)
	if err != nil {
		panic(err)
	}

	var allShortenedUrls []string
	for _, shortenedUrlsAsString := range shortenedUrlsCreatedByUsers {
		var shortenedUrls []string
		json.Unmarshal([]byte(shortenedUrlsAsString), &shortenedUrls)
		allShortenedUrls = append(allShortenedUrls, shortenedUrls...)
	}

	allKeys, err := r.urlRepo.Keys(c, "*")
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(allKeys); i++ {
		allKeys[i] = strings.TrimPrefix(allKeys[i], viper.GetString("database.redis.prefix")+"_")
	}
	// helps us get the keys that have ttl>0
	allShortenedUrls = findCommonElementsInSlices(allShortenedUrls, allKeys)

	return allShortenedUrls, nil
}

func findCommonElementsInSlices(slice1, slice2 []string) []string {
	elementsMap := make(map[string]bool)

	for _, elem := range slice1 {
		elementsMap[elem] = true
	}

	var commonElements []string

	for _, elem := range slice2 {
		// If the element exists in the map, it's common
		if elementsMap[elem] {
			commonElements = append(commonElements, elem)
		}
	}

	return commonElements
}

func (r *urlRepository) GetAllShortenedUrlKeysForUser(c context.Context, userId string) ([]string, error) {
	userId = "User" + userId
	shortenedUrlsCreatedByUser, err := r.urlRepo.HGet(c, userId, "urlsList")
	if err != nil {
		panic(err)
	}
	var shortenedUrlsList []string
	if shortenedUrlsCreatedByUser == "" {
		return shortenedUrlsList, nil
	}
	err = json.Unmarshal([]byte(shortenedUrlsCreatedByUser), &shortenedUrlsList)
	if err != nil {
		panic(err)
	}

	allKeys, err := r.urlRepo.Keys(c, "*")
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(allKeys); i++ {
		allKeys[i] = strings.TrimPrefix(allKeys[i], viper.GetString("database.redis.prefix")+"_")
	}
	// helps us get the keys that have ttl>0
	shortenedUrlsList = findCommonElementsInSlices(shortenedUrlsList, allKeys)

	return shortenedUrlsList, nil
}

func (r *urlRepository) GetUrlForShortenedKeys(c context.Context, shortenedKeys []string) (map[string]string, error) {
	var redisKeysWithField = make(map[string]string)
	for i := 0; i < len(shortenedKeys); i++ {
		redisKeysWithField[shortenedKeys[i]] = "url"
	}

	return r.urlRepo.HGets(c, redisKeysWithField)
}

func (r *urlRepository) GetTtlForShortenedKeys(c context.Context, shortenedKeys []string) (map[string]uint64, error) {
	var redisKeysWithField = make(map[string]string)
	for i := 0; i < len(shortenedKeys); i++ {
		redisKeysWithField[shortenedKeys[i]] = "toBeDeletedAt"
	}

	shoretenedKeyWithDeletedAt, err := r.urlRepo.HGets(c, redisKeysWithField)
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

	shoretenedKeyWithUrlHits, err := r.urlRepo.HGets(c, redisKeysWithField)
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

	return r.urlRepo.HGets(c, redisKeysWithField)
}

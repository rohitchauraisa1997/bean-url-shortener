package services

import (
	"context"
	"strings"
	"time"

	// bean "github.com/retail-ai-inc/bean/v2/trace"
	"backend/helpers"
	"backend/models"
	"backend/repositories"

	"github.com/spf13/viper"
)

type UrlService interface {
	ShortenUrl(ctx context.Context, urlToShorten string, userIp string, expiryDuration time.Duration) (string, error)
	ResolveUrl(ctx context.Context, shortenedUrl string) (string, error)
	// GetAnalyticsForAllUsers(ctx context.Context) (models.ShortenedUrlAndDetailsSlice, error)
	GetAnalyticsForAllUsers(ctx context.Context) (map[string]models.ShortenedUrlAndDetailsSlice, error)
	GetAnalyticsForUser(ctx context.Context, userId string) (models.ShortenedUrlAndDetailsSlice, error)
}

type urlService struct {
	urlRepository repositories.UrlRepository
}

func NewUrlService(urlRepo repositories.UrlRepository) *urlService {
	return &urlService{
		urlRepository: urlRepo,
	}
}

func (service *urlService) ShortenUrl(ctx context.Context, urlToShorten string, userIp string, expiryDurationInMins time.Duration) (string, error) {
	shortenedUrl, err := service.urlRepository.GetNewShortenedUrl(ctx)
	if err != nil {
		panic(err)
	}
	service.urlRepository.CreateShortenedUrl(ctx, shortenedUrl, urlToShorten, userIp, expiryDurationInMins)
	return shortenedUrl, nil
}

func (service *urlService) ResolveUrl(ctx context.Context, shortenedUrl string) (string, error) {
	url, err := service.urlRepository.GetUrlViaShortenedUrl(ctx, shortenedUrl)
	if err != nil {
		panic(err)
	}
	service.urlRepository.IncreaseVisitCounter(ctx, shortenedUrl)
	return url, nil
}

func (service *urlService) GetAnalyticsForAllUsers(ctx context.Context) (map[string]models.ShortenedUrlAndDetailsSlice, error) {
	shortenedKeys, _ := service.urlRepository.GetAllShortenedUrlKeysForAllUsers(ctx)
	shortenedKeysWithUrls, _ := service.urlRepository.GetUrlForShortenedKeys(ctx, shortenedKeys)
	shortenedKeysWithTtls, _ := service.urlRepository.GetTtlForShortenedKeys(ctx, shortenedKeys)
	shortenedKeysWithUrlHits, _ := service.urlRepository.GetUrlHitsForShortenedKeys(ctx, shortenedKeys)
	shortenedKeysWithCreatedBy, _ := service.urlRepository.GetCreatedByForShortenedKeys(ctx, shortenedKeys)

	domainName := viper.GetString("backend.domain")
	keyPrefix := viper.GetString("database.redis.prefix")

	var adminResponse = make(map[string]models.ShortenedUrlAndDetailsSlice)
	for shortenedKey, createdByUser := range shortenedKeysWithCreatedBy {
		var shortenedUrlAnalytics models.UrlAnalyticDetails
		if val, ok := shortenedKeysWithUrls[shortenedKey]; ok {
			val = strings.Trim(val, "\"")
			shortenedUrlAnalytics.URL = val
		}
		if val, ok := shortenedKeysWithTtls[shortenedKey]; ok {
			ttl := uint64(int64(val) - time.Now().Unix())
			shortenedUrlAnalytics.TTL = ttl
		}
		if val, ok := shortenedKeysWithUrlHits[shortenedKey]; ok {
			shortenedUrlAnalytics.Hits = val
		}
		if val, ok := shortenedKeysWithCreatedBy[shortenedKey]; ok {
			val = strings.Trim(val, "\"")
			shortenedUrlAnalytics.CreatedBy = val
		}

		// if createdByUser already in adminResponse add the shortenedUrl to the previously added slice
		if _, ok := adminResponse[createdByUser]; ok {
			shortenedKey = strings.TrimPrefix(shortenedKey, keyPrefix+"_")
			shortenedKey = domainName + "/" + shortenedKey
			adminResponse[createdByUser] = append(adminResponse[createdByUser], models.ShortenedUrlAndDetail{
				ShortenedUrl:  shortenedKey,
				UrlsAnalytics: shortenedUrlAnalytics,
			})
		} else {
			// initiate the slice
			shortenedKey = strings.TrimPrefix(shortenedKey, keyPrefix+"_")
			shortenedKey = domainName + "/" + shortenedKey
			adminResponse[createdByUser] = models.ShortenedUrlAndDetailsSlice{models.ShortenedUrlAndDetail{
				ShortenedUrl:  shortenedKey,
				UrlsAnalytics: shortenedUrlAnalytics,
			}}
		}
	}

	return adminResponse, nil
}

func (service *urlService) GetAnalyticsForUser(ctx context.Context, userId string) (models.ShortenedUrlAndDetailsSlice, error) {
	shortenedKeys, _ := service.urlRepository.GetAllShortenedUrlKeysForUser(ctx, userId)

	shortenedKeysWithUrls, _ := service.urlRepository.GetUrlForShortenedKeys(ctx, shortenedKeys)
	shortenedKeysWithTtls, _ := service.urlRepository.GetTtlForShortenedKeys(ctx, shortenedKeys)
	shortenedKeysWithUrlHits, _ := service.urlRepository.GetUrlHitsForShortenedKeys(ctx, shortenedKeys)
	shortenedKeysWithCreatedBy, _ := service.urlRepository.GetCreatedByForShortenedKeys(ctx, shortenedKeys)

	var fulResponse = make(map[string]models.UrlAnalyticDetails)
	domainName := viper.GetString("backend.domain")
	keyPrefix := viper.GetString("database.redis.prefix") + "_"

	for _, skey := range shortenedKeys {
		prefixedSkey := keyPrefix + skey
		var shortenedUrlAnalytics models.UrlAnalyticDetails
		if val, ok := shortenedKeysWithUrls[prefixedSkey]; ok {
			val = strings.Trim(val, "\"")
			shortenedUrlAnalytics.URL = val
		}
		if val, ok := shortenedKeysWithTtls[prefixedSkey]; ok {
			ttl := uint64(int64(val) - time.Now().Unix())
			shortenedUrlAnalytics.TTL = ttl
		}
		if val, ok := shortenedKeysWithUrlHits[prefixedSkey]; ok {
			shortenedUrlAnalytics.Hits = val
		}
		if val, ok := shortenedKeysWithCreatedBy[prefixedSkey]; ok {
			val = strings.Trim(val, "\"")
			shortenedUrlAnalytics.CreatedBy = val
		}
		skey = domainName + "/" + skey
		fulResponse[skey] = shortenedUrlAnalytics
	}
	sortedResponse := helpers.SortResponse(fulResponse)
	return sortedResponse, nil
}

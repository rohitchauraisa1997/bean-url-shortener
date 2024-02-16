package services

import (
	"context"
	"encoding/json"
	"fmt"
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
	GetAnalytics(ctx context.Context) (models.ShortenedUrlAndDetailsSlice, error)
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

func (service *urlService) GetAnalytics(ctx context.Context) (models.ShortenedUrlAndDetailsSlice, error) {
	shortenedKeys, _ := service.urlRepository.GetAllShortenedUrlKeysForAllUsers(ctx)
	jsonshortenedKeys, _ := json.MarshalIndent(shortenedKeys, " ", " ")
	fmt.Println("*************")
	fmt.Println(string(jsonshortenedKeys))
	fmt.Println("*************")
	shortenedKeysWithUrls, _ := service.urlRepository.GetUrlForShortenedKeys(ctx, shortenedKeys)
	shortenedKeysWithTtls, _ := service.urlRepository.GetTtlForShortenedKeys(ctx, shortenedKeys)
	shortenedKeysWithUrlHits, _ := service.urlRepository.GetUrlHitsForShortenedKeys(ctx, shortenedKeys)
	shortenedKeysWithCreatedBy, _ := service.urlRepository.GetCreatedByForShortenedKeys(ctx, shortenedKeys)

	var fulResponse = make(map[string]models.UrlAnalyticDetails)
	domainName := viper.GetString("backend.domain")

	for _, skey := range shortenedKeys {
		prefixedSkey := viper.GetString("database.redis.prefix") + "_" + skey
		var routeResolutionDetails models.UrlAnalyticDetails
		if val, ok := shortenedKeysWithUrls[prefixedSkey]; ok {
			val = strings.Trim(val, "\"")
			routeResolutionDetails.URL = val
		}
		if val, ok := shortenedKeysWithTtls[prefixedSkey]; ok {
			ttl := uint64(int64(val) - time.Now().Unix())
			routeResolutionDetails.TTL = ttl
		}
		if val, ok := shortenedKeysWithUrlHits[prefixedSkey]; ok {
			routeResolutionDetails.Hits = val
		}
		if val, ok := shortenedKeysWithCreatedBy[prefixedSkey]; ok {
			val = strings.Trim(val, "\"")
			routeResolutionDetails.CreatedBy = val
		}
		skey = domainName + "/" + skey
		fulResponse[skey] = routeResolutionDetails
	}
	sortedResponse := helpers.SortResponse(fulResponse)
	return sortedResponse, nil
}

func (service *urlService) GetAnalyticsForUser(ctx context.Context, userId string) (models.ShortenedUrlAndDetailsSlice, error) {
	shortenedKeys, _ := service.urlRepository.GetAllShortenedUrlKeysForUser(ctx, userId)

	shortenedKeysWithUrls, _ := service.urlRepository.GetUrlForShortenedKeys(ctx, shortenedKeys)
	shortenedKeysWithTtls, _ := service.urlRepository.GetTtlForShortenedKeys(ctx, shortenedKeys)
	shortenedKeysWithUrlHits, _ := service.urlRepository.GetUrlHitsForShortenedKeys(ctx, shortenedKeys)
	shortenedKeysWithCreatedBy, _ := service.urlRepository.GetCreatedByForShortenedKeys(ctx, shortenedKeys)

	var fulResponse = make(map[string]models.UrlAnalyticDetails)
	domainName := viper.GetString("backend.domain")

	for _, skey := range shortenedKeys {
		prefixedSkey := viper.GetString("database.redis.prefix") + "_" + skey
		var routeResolutionDetails models.UrlAnalyticDetails
		if val, ok := shortenedKeysWithUrls[prefixedSkey]; ok {
			val = strings.Trim(val, "\"")
			// val = strings.TrimPrefix(val, viper.GetString("database.redis.prefix")+"_")
			routeResolutionDetails.URL = val
		}
		if val, ok := shortenedKeysWithTtls[prefixedSkey]; ok {
			ttl := uint64(int64(val) - time.Now().Unix())
			routeResolutionDetails.TTL = ttl
		}
		if val, ok := shortenedKeysWithUrlHits[prefixedSkey]; ok {
			routeResolutionDetails.Hits = val
		}
		if val, ok := shortenedKeysWithCreatedBy[prefixedSkey]; ok {
			val = strings.Trim(val, "\"")
			routeResolutionDetails.CreatedBy = val
		}
		skey = domainName + "/" + skey
		fulResponse[skey] = routeResolutionDetails
	}
	sortedResponse := helpers.SortResponse(fulResponse)
	return sortedResponse, nil
}

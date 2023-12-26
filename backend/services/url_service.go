package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	// "github.com/retail-ai-inc/bean/trace"
	"backend/helpers"
	"backend/models"
	"backend/repositories"

	"github.com/spf13/viper"
)

type UrlService interface {
	ShortenUrl(ctx context.Context, urlToShorten string, userIp string, expiryDuration time.Duration) (string, error)
	ResolveUrl(ctx context.Context, shortenedUrl string) (string, error)
	GetAnalytics(ctx context.Context) (models.ShortenedUrlAndDetailsSlice, error)
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
	shortenedKeys, _ := service.urlRepository.GetAllShortenedUrlKeys(ctx)
	shortenedKeysWithUrls, _ := service.urlRepository.GetUrlForShortenedKeys(ctx, shortenedKeys)
	shortenedKeysWithTtls, _ := service.urlRepository.GetTtlForShortenedKeys(ctx, shortenedKeys)
	shortenedKeysWithUrlHits, _ := service.urlRepository.GetUrlHitsForShortenedKeys(ctx, shortenedKeys)
	shortenedKeysWithCreatedBy, _ := service.urlRepository.GetCreatedByForShortenedKeys(ctx, shortenedKeys)

	var fulResponse = make(map[string]models.UrlAnalyticDetails)
	domainName := viper.GetString("backend.domain")

	for _, skey := range shortenedKeys {
		var routeResolutionDetails models.UrlAnalyticDetails
		if val, ok := shortenedKeysWithUrls[skey]; ok {
			fmt.Println("val", val)
			val = strings.Trim(val, "\"")
			routeResolutionDetails.URL = val
		}
		if val, ok := shortenedKeysWithTtls[skey]; ok {
			fmt.Println("ttlval", val)
			ttl := uint64(int64(val) - time.Now().Unix())
			routeResolutionDetails.TTL = ttl
		}
		if val, ok := shortenedKeysWithUrlHits[skey]; ok {
			fmt.Println("urlHitsVal", val)
			routeResolutionDetails.Hits = val
		}
		if val, ok := shortenedKeysWithCreatedBy[skey]; ok {
			fmt.Println("createdByval", val)
			val = strings.Trim(val, "\"")
			routeResolutionDetails.CreatedBy = val
		}
		skey = domainName + "/" + skey
		fulResponse[skey] = routeResolutionDetails
	}
	sortedResponse := helpers.SortResponse(fulResponse)
	return sortedResponse, nil
}

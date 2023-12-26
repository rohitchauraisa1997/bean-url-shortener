package helpers

import (
	"backend/models"
	"sort"
)

func SortResponse(resp map[string]models.UrlAnalyticDetails) models.ShortenedUrlAndDetailsSlice {
	// Convert map to a slice for sorting
	var items models.ShortenedUrlAndDetailsSlice
	for key, value := range resp {
		items = append(items, models.ShortenedUrlAndDetail{ShortenedUrl: key, UrlsAnalytics: value})
	}

	// Sort the slice in descending order based on Hits
	sort.Sort(items)

	return items
}

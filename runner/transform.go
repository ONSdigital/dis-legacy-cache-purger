package runner

import (
	"context"
	"fmt"
	"strings"

	"github.com/ONSdigital/dp-legacy-cache-api/models"
)

type CollectionCachePaths struct {
	collectionID    string
	collectionTitle string
	paths           []string
}

func mapCacheTimeByCollection(ctx context.Context, cacheTimes []*models.CacheTime, domains []string) map[string]CollectionCachePaths {
	result := make(map[string]CollectionCachePaths)
	for _, ct := range cacheTimes {
		if ct == nil {
			continue
		}
		collectionID := ct.CollectionID
		for _, domain := range domains {
			prefixedPath := fmt.Sprintf("%s%s", domain, ct.Path)
			collectionCacheTimes := result[collectionID]
			collectionCacheTimes.collectionID = collectionID
			collectionCacheTimes.collectionTitle = ct.CollectionTitle
			collectionCacheTimes.paths = append(collectionCacheTimes.paths, prefixedPath)
			result[collectionID] = collectionCacheTimes
		}
	}
	return result
}

func mapCollectionCacheTimeMapToRequests(ctx context.Context, collectionCacheTimeMap map[string]CollectionCachePaths) []CollectionCachePurgeRequest {
	requests := make([]CollectionCachePurgeRequest, 0, len(collectionCacheTimeMap))
	for collectionID, collectionCachePaths := range collectionCacheTimeMap {
		var prefixes []string
		var files []string
		for _, path := range collectionCachePaths.paths {
			if strings.Contains(path, "/timeseries/") {
				// exclude timeseries paths.
				continue
			}

			// Add standard path.
			files = append(files, fmt.Sprintf("https://%s", path))

			// If the path does not contain a query string, we can also purge the /data and /pdf versions of the file.
			if !strings.Contains(path, "?") {
				files = append(files,
					fmt.Sprintf("https://%s/data", path),
					fmt.Sprintf("https://%s/pdf", path),
				)
			}
		}
		requests = append(requests, CollectionCachePurgeRequest{
			CollectionID:    collectionID,
			CollectionTitle: collectionCachePaths.collectionTitle,
			Prefixes:        prefixes,
			Files:           files,
		})
	}
	return requests
}

func transformCacheTimesToCollectionCachePurgeRequests(ctx context.Context, cacheTimes []*models.CacheTime, domains []string) []CollectionCachePurgeRequest {
	cacheTimesByCollectionID := mapCacheTimeByCollection(ctx, cacheTimes, domains)
	return mapCollectionCacheTimeMapToRequests(ctx, cacheTimesByCollectionID)
}

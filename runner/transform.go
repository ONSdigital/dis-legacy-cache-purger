package runner

import (
	"context"
	"fmt"
	"strings"

	"github.com/ONSdigital/dp-legacy-cache-api/models"
)

func mapCacheTimeByCollectionID(ctx context.Context, cacheTimes []*models.CacheTime, domains []string) map[string][]string {
	result := make(map[string][]string)
	for _, ct := range cacheTimes {
		if ct == nil {
			continue
		}
		collectionID := ct.CollectionID
		for _, domain := range domains {
			prefixedPath := fmt.Sprintf("%s%s", domain, ct.Path)
			result[collectionID] = append(result[collectionID], prefixedPath)
		}
	}
	return result
}

func mapCollectionCacheTimeMapToRequests(ctx context.Context, collectionCacheTimeMap map[string][]string) []CollectionCachePurgeRequest {
	requests := make([]CollectionCachePurgeRequest, 0, len(collectionCacheTimeMap))
	for collectionID, paths := range collectionCacheTimeMap {
		var prefixes []string
		var files []string
		for _, path := range paths {
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
			CollectionID: collectionID,
			Prefixes:     prefixes,
			Files:        files,
		})
	}
	return requests
}

func transformCacheTimesToCollectionCachePurgeRequests(ctx context.Context, cacheTimes []*models.CacheTime, domains []string) []CollectionCachePurgeRequest {
	cacheTimesByCollectionID := mapCacheTimeByCollectionID(ctx, cacheTimes, domains)
	return mapCollectionCacheTimeMapToRequests(ctx, cacheTimesByCollectionID)
}

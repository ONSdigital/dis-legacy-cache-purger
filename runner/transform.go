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

func pathMatches(path string, filePaths []string) bool {
	for _, file := range filePaths {
		if strings.Contains(path, file) {
			return true
		}
	}
	return false
}

func pathToFiles(path string, dataPaths, pdfPaths, excludedPaths []string) []string {
	if pathMatches(path, excludedPaths) {
		return nil
	}
	files := []string{fmt.Sprintf("https://%s", path)}
	if !strings.Contains(path, "?") {
		if pathMatches(path, dataPaths) {
			files = append(files, fmt.Sprintf("https://%s/data", path))
		}
		if pathMatches(path, pdfPaths) {
			files = append(files, fmt.Sprintf("https://%s/pdf", path))
		}
	}
	return files
}

func mapCollectionCacheTimeMapToRequests(ctx context.Context, collectionCacheTimeMap map[string]CollectionCachePaths, dataPaths, pdfPaths, excludedPaths []string) []CollectionCachePurgeRequest {
	requests := make([]CollectionCachePurgeRequest, 0, len(collectionCacheTimeMap))
	for collectionID, collectionCachePaths := range collectionCacheTimeMap {
		var prefixes []string
		files := make([]string, 0, len(collectionCachePaths.paths))
		for _, path := range collectionCachePaths.paths {
			files = append(files, pathToFiles(path, dataPaths, pdfPaths, excludedPaths)...)
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

func transformCacheTimesToCollectionCachePurgeRequests(ctx context.Context, cacheTimes []*models.CacheTime, domains, dataPaths, pdfPaths, excludedPaths []string) []CollectionCachePurgeRequest {
	cacheTimesByCollectionID := mapCacheTimeByCollection(ctx, cacheTimes, domains)
	return mapCollectionCacheTimeMapToRequests(ctx, cacheTimesByCollectionID, dataPaths, pdfPaths, excludedPaths)
}

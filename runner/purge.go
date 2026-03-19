package runner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/cache"
)

func (p *PurgeRunner) CachePurge(ctx context.Context, requests []CollectionCachePurgeRequest, releaseTime time.Time) []CachePurgeResult {
	cachePurgeResults := make([]CachePurgeResult, 0, len(requests))

	for _, req := range requests {
		cachePurgeResult := p.CachePurgeCollection(ctx, req, releaseTime)
		cachePurgeResults = append(cachePurgeResults, cachePurgeResult)
	}

	return cachePurgeResults
}

func (p *PurgeRunner) CachePurgeCollection(ctx context.Context, req CollectionCachePurgeRequest, releaseTime time.Time) CachePurgeResult {
	logData := log.Data{"collection_id": req.CollectionID}

	log.Info(ctx, "purging cache for collection", logData)

	var slackErr error

	successfulPathsPurged := 0

	prefixErr := p.CachePurgePrefixes(ctx, req.Prefixes)
	if prefixErr == nil {
		successfulPathsPurged += len(req.Prefixes)
	}

	fileErr := p.CachePurgeFiles(ctx, req.Files)
	if fileErr == nil {
		successfulPathsPurged += len(req.Files)
	}

	if prefixErr != nil || fileErr != nil {
		log.Error(ctx, "error purging cache for collection", fileErr, logData)
		slackErr = p.sendFailureMessageForCollection(ctx, req, releaseTime)
		if slackErr != nil {
			log.Error(ctx, "error sending Slack message after cache purge failure", slackErr, logData)
		}

		return CachePurgeResult{
			Purges:  successfulPathsPurged,
			Success: false,
			Error:   fmt.Errorf("errors occurred during cache purge"),
		}
	} else {
		slackErr = p.sendSuccessMessageForCollection(ctx, req, releaseTime)
		if slackErr != nil {
			log.Error(ctx, "error sending Slack message after cache purge failure", slackErr, logData)
		}

		log.Info(ctx, "completed purging cache for collection", logData)

		return CachePurgeResult{
			Purges:  successfulPathsPurged,
			Success: true,
			Error:   nil,
		}
	}
}

func (p *PurgeRunner) CachePurgePrefixes(ctx context.Context, prefixes []string) error {
	batchSize := p.config.CloudflareBatchSize
	maxParallel := p.config.MaxParallel
	return batchProcess(ctx, prefixes, batchSize, maxParallel, func(ctx context.Context, batch []string) error {
		prefixPurgeParams := cache.CachePurgeParams{
			ZoneID: cloudflare.F(p.config.CloudflareZoneID),
			Body: cache.CachePurgeParamsBodyCachePurgeFlexPurgeByPrefixes{
				Prefixes: cloudflare.F(batch),
			},
		}
		_, err := p.clientList.CloudflareCacheClient.Purge(ctx, prefixPurgeParams)
		return err
	})
}

func (p *PurgeRunner) CachePurgeFiles(ctx context.Context, files []string) error {
	batchSize := p.config.CloudflareBatchSize
	maxParallel := p.config.MaxParallel
	return batchProcess(ctx, files, batchSize, maxParallel, func(ctx context.Context, batch []string) error {
		filePurgeParams := cache.CachePurgeParams{
			ZoneID: cloudflare.F(p.config.CloudflareZoneID),
			Body: cache.CachePurgeParamsBodyCachePurgeSingleFile{
				Files: cloudflare.F(batch),
			},
		}
		_, err := p.clientList.CloudflareCacheClient.Purge(ctx, filePurgeParams)
		return err
	})
}

// batchProcess splits a slice of items into batches of the given size and processes them concurrently with a maxParallel limit.
func batchProcess[T any](ctx context.Context, items []T, batchSize, maxParallel int, fn func(ctx context.Context, batch []T) error) error {
	if maxParallel < 1 {
		maxParallel = 1
	}

	if batchSize < 1 {
		batchSize = 1
	}

	sem := make(chan struct{}, maxParallel)
	var wg sync.WaitGroup
	var firstErr error
	var errOnce sync.Once

	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}
		batch := items[i:end]

		sem <- struct{}{}
		wg.Add(1)
		go func(batch []T) {
			defer wg.Done()
			defer func() { <-sem }()
			if err := fn(ctx, batch); err != nil {
				errOnce.Do(func() { firstErr = err })
			}
		}(batch)
	}

	wg.Wait()
	return firstErr
}

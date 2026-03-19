package runner

import (
	"context"
	"sync"
	"time"

	"github.com/ONSdigital/dp-legacy-cache-api/models"
	"github.com/ONSdigital/dp-legacy-cache-api/sdk"
	"github.com/ONSdigital/log.go/v2/log"
)

const (
	//TODO: consider making this configurable if needed
	legacyCacheAPIBatchSize = 100
)

func (p *PurgeRunner) GetCacheTimes(ctx context.Context, releaseTime time.Time) ([]*models.CacheTime, error) {
	auth := sdk.Auth{
		ServiceAuthToken: p.config.ServiceToken,
	}

	// Fetch the first page to get total count and items
	firstOpts := sdk.Options{
		ReleaseTime: releaseTime,
		Offset:      0,
		Limit:       legacyCacheAPIBatchSize,
	}

	firstPage, err := p.clientList.LegacyCacheClient.GetCacheTimes(ctx, auth, firstOpts)
	if err != nil {
		log.Error(ctx, "error getting cache times", err)
		return []*models.CacheTime{}, err
	}

	totalCacheTimes := firstPage.TotalCount
	cacheTimes := make([]*models.CacheTime, 0, totalCacheTimes)
	cacheTimes = append(cacheTimes, firstPage.Items...)

	if totalCacheTimes <= len(firstPage.Items) {
		return cacheTimes, nil
	}

	// Parallelize fetching remaining pages using goroutines and a wait group, with configurable max parallelism
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		errOnce  sync.Once
		fetchErr error
	)

	maxParallel := p.config.MaxParallel
	if maxParallel < 1 {
		maxParallel = 1
	}

	sem := make(chan struct{}, maxParallel)

	// Start from offset = len(firstPage.Items)
	for offset := len(firstPage.Items); offset < totalCacheTimes; offset += legacyCacheAPIBatchSize {
		sem <- struct{}{} // acquire a slot
		wg.Add(1)
		go func(offset int) {
			defer wg.Done()
			defer func() { <-sem }() // release the slot
			opts := sdk.Options{
				Offset:      offset,
				Limit:       legacyCacheAPIBatchSize,
				ReleaseTime: releaseTime,
			}
			page, err := p.clientList.LegacyCacheClient.GetCacheTimes(ctx, auth, opts)
			if err != nil {
				errOnce.Do(func() {
					fetchErr = err
				})
				return
			}
			mu.Lock()
			cacheTimes = append(cacheTimes, page.Items...)
			mu.Unlock()
		}(offset)
	}

	wg.Wait()
	return cacheTimes, fetchErr
}

package runner

import (
	"context"
	"time"

	"github.com/ONSdigital/log.go/v2/log"
)

// Run is the main logic of the app. It gets cachetimes scheduled for publish and then attempts to purge them.
func (p *PurgeRunner) Run(ctx context.Context, publishTime time.Time) (*PurgeRunResult, error) {
	// The time to check for scheduled publication, this is rounded to the nearest minute as publication on the minute
	// is what is provided to users to enter.  Validation is carried out below to ensure publications are not made early

	logData := log.Data{"release_time": publishTime}

	log.Info(ctx, "retrieving cache times scheduled for publish", logData)

	cacheTimes, fetchErr := p.GetCacheTimes(ctx, publishTime)
	if fetchErr != nil {
		log.Error(ctx, "error getting cache times", fetchErr, logData)
		return &PurgeRunResult{}, fetchErr
	}

	totalCacheTimes := len(cacheTimes)
	if totalCacheTimes == 0 {
		log.Info(ctx, "no cache times ready for purging")
		return &PurgeRunResult{
			Success: true,
			Results: []CachePurgeResult{},
		}, nil
	}

	cachePurgeRequests := transformCacheTimesToCollectionCachePurgeRequests(ctx, cacheTimes, p.config.Domains)

	for _, req := range cachePurgeRequests {
		log.Info(ctx, "waiting to purge cache for collection", log.Data{
			"collection_id": req.CollectionID,
			"total_paths":   len(req.Prefixes) + len(req.Files),
			"prefixes":      len(req.Prefixes),
			"files":         len(req.Files),
		})
	}

	// Get the time difference between the minute submitted in the query and the current time as specified above
	purgeCheck := publishTime.Sub(time.Now().Add(-p.config.CachePurgeDiffTime))

	// Check to ensure cache is not purged early - sleeps the process for the amount of time between current time
	// above and publication time
	p.config.SleepFunc(purgeCheck)
	results := p.CachePurge(ctx, cachePurgeRequests, publishTime)

	log.Info(ctx, "completed cache purge run")

	return &PurgeRunResult{
		Success: true,
		Results: results,
	}, nil
}

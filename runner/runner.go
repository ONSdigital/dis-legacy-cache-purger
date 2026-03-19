package runner

import (
	"context"
	"time"

	"github.com/ONSdigital/dis-legacy-cache-purger/clients"
	config "github.com/ONSdigital/dis-legacy-cache-purger/config"
)

type CachePurgeResult struct {
	Purges  int
	Success bool
	Error   error
}

type PurgeRunResult struct {
	Results []CachePurgeResult
	Success bool
}

type CollectionCachePurgeRequest struct {
	CollectionID string
	Prefixes     []string
	Files        []string
}

type Runner interface {
	Run(ctx context.Context, publishTime time.Time) (*PurgeRunResult, error)
}

// Runner is the main logic/orchestrator of the application.
type PurgeRunner struct {
	clientList clients.ClientList
	config     *config.Configuration
}

func NewPurgeRunner(cfg *config.Configuration, clientList clients.ClientList) (Runner, error) {
	return &PurgeRunner{
		clientList: clientList,
		config:     cfg,
	}, nil
}

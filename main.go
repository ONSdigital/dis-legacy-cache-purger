package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ONSdigital/dis-legacy-cache-purger/clients"
	mockClients "github.com/ONSdigital/dis-legacy-cache-purger/clients/mock"
	"github.com/ONSdigital/dis-legacy-cache-purger/config"
	"github.com/ONSdigital/dis-legacy-cache-purger/runner"
	"github.com/ONSdigital/dp-legacy-cache-api/models"
	"github.com/ONSdigital/dp-legacy-cache-api/sdk"
	apiError "github.com/ONSdigital/dp-legacy-cache-api/sdk/errors"
	sdkMock "github.com/ONSdigital/dp-legacy-cache-api/sdk/mocks"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/cache"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

const serviceName = "dis-legacy-cache-purger"

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

func main() {
	log.Namespace = serviceName
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(ctx, "fatal runtime error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return errors.Wrap(err, "unable to retrieve service configuration")
	}
	log.Info(ctx, "config on startup", log.Data{"config": cfg, "build_time": BuildTime, "git-commit": GitCommit})

	clList, err := newClientListFromConfig(cfg)
	if err != nil {
		return errors.Wrap(err, "unable to create client list from config")
	}
	purgeRunner, err := runner.NewPurgeRunner(cfg, *clList)

	if err != nil {
		return errors.Wrap(err, "unable to instantiate purge runner")
	}

	// Run the runner in the background, using a result channel and an error channel for fatal errors
	errChan := make(chan error, 1)
	resultChan := make(chan *runner.PurgeRunResult, 1)
	go func() {
		now := time.Now().UTC()
		// get difference between now and next minute - to the whole minute to improve accuracy
		nextMinute := now.Truncate(time.Minute).Add(time.Minute)

		result, runErr := purgeRunner.Run(ctx, nextMinute)
		if runErr != nil {
			errChan <- runErr
		}
		resultChan <- result
	}()

	// blocks until completion, an os interrupt or a fatal error occurs
	select {
	case err = <-errChan:
		log.Error(ctx, "runner error received", err)
		return err
	case sig := <-signals:
		log.Info(ctx, "os signal received", log.Data{"signal": sig})
	case result := <-resultChan:
		log.Info(ctx, "purge cache for schedule result", log.Data{"result": result})
		if !result.Success {
			if err != nil {
				log.Error(ctx, "unable to send notification of result", err)
				return err
			}
		}
		log.Info(ctx, "purge cache for schedule complete")
	}
	return nil // TODO close down the checker and confirm task completion state (err or nil)
}

func newClientListFromConfig(cfg *config.Configuration) (*clients.ClientList, error) {
	cloudflareClient := getCloudflareClient(cfg.EnableCloudflarePurge, cfg.CloudflareAPIToken)
	legacyCacheAPIClient := getLegacyCacheAPIClient(cfg.EnableCacheAPI, cfg.LegacyCacheAPIURL)
	slackClient := getSlackClient(cfg.EnableSlackAlerts, cfg.SlackAPIToken)

	return &clients.ClientList{
		CloudflareCacheClient: cloudflareClient,
		LegacyCacheClient:     legacyCacheAPIClient,
		SlackClient:           slackClient,
	}, nil
}

func getCloudflareClient(enabled bool, token string) clients.CloudflareCacheClienter {
	var cloudflareClient clients.CloudflareCacheClienter

	if enabled {
		cloudflareClient = cloudflare.NewClient(
			option.WithAPIToken(token),
		).Cache
	} else {
		cloudflareClient = &mockClients.CloudflareCacheClienterMock{
			PurgeFunc: func(ctx context.Context, params cache.CachePurgeParams, opts ...option.RequestOption) (*cache.CachePurgeResponse, error) {
				log.Info(ctx, "mock Cloudflare cache purge called", log.Data{"params": params, "opts": opts})
				return &cache.CachePurgeResponse{}, nil
			},
		}
	}
	return cloudflareClient
}

func getLegacyCacheAPIClient(enabled bool, url string) sdk.Clienter {
	var legacyCacheAPIClient sdk.Clienter

	if enabled {
		legacyCacheAPIClient = sdk.New(url)
	} else {
		legacyCacheAPIClient = &sdkMock.ClienterMock{
			GetCacheTimesFunc: func(ctx context.Context, auth sdk.Auth, options sdk.Options) (*models.CacheTimesList, apiError.Error) {
				log.Info(ctx, "mock Legacy Cache API get cache times called", log.Data{"options": options})
				releaseTime := time.Now().Add(1 * time.Minute)
				return &models.CacheTimesList{
					Count:      1,
					TotalCount: 1,
					Items:      generateTestData(3, 1, releaseTime),
				}, nil
			},
		}
	}
	return legacyCacheAPIClient
}

func getSlackClient(enabled bool, token string) clients.SlackClienter {
	var slackClient clients.SlackClienter

	if enabled {
		slackClient = slack.New(token)
	} else {
		slackClient = &mockClients.SlackClienterMock{
			PostMessageFunc: func(channel string, options ...slack.MsgOption) (string, string, error) {
				log.Info(context.Background(), "mock Slack post message called", log.Data{"channel": channel})
				return "mockTimestamp", channel, nil
			},
		}
	}
	return slackClient
}

func generateTestData(numTimes, numCollections int, releaseTime time.Time) []*models.CacheTime {
	testData := make([]*models.CacheTime, 0, numTimes*numCollections)
	rnd := time.Now().UnixNano()
	for c := 1; c <= numCollections; c++ {
		collectionID := fmt.Sprintf("test-collection-%d", c)
		for t := 1; t <= numTimes; t++ {
			id := fmt.Sprintf("%d-%d", c, t)
			// Randomly generate a path
			base := []string{"/economy", "/health", "/population", "/business", "/education", "/housing"}
			basePath := base[(int(rnd)+c+t)%len(base)]
			path := fmt.Sprintf("%s/%d", basePath, (int(rnd)+c*t)%100)
			// 40% chance to add a query string
			if ((int(rnd) + c*t + t) % 10) < 4 {
				path = fmt.Sprintf("%s?%s=%d", path, randomQueryKey((int(rnd)+c*t)%5), (int(rnd)+c*t+t)%1000)
			}
			testData = append(testData, &models.CacheTime{
				ID:           id,
				CollectionID: collectionID,
				Path:         path,
				ReleaseTime:  &releaseTime,
			})
		}
	}
	return testData
}

func randomQueryKey(seed int) string {
	keys := []string{"q", "search", "id", "type", "filter"}
	return keys[seed%len(keys)]
}

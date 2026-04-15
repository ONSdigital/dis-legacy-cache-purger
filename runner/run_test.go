package runner

import (
	"context"
	"testing"
	"time"

	"github.com/ONSdigital/dis-legacy-cache-purger/clients"
	mockClients "github.com/ONSdigital/dis-legacy-cache-purger/clients/mock"
	apiError "github.com/ONSdigital/dp-legacy-cache-api/sdk/errors"
	"github.com/cloudflare/cloudflare-go/v6/cache"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/slack-go/slack"

	"github.com/ONSdigital/dp-legacy-cache-api/sdk"
	sdkMock "github.com/ONSdigital/dp-legacy-cache-api/sdk/mocks"

	config "github.com/ONSdigital/dis-legacy-cache-purger/config"
	"github.com/ONSdigital/dp-legacy-cache-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRunnerRun(t *testing.T) {
	Convey("Given a Runner with a page to purge the cache for and mock clients that don't error", t, func() {
		publishTime := time.Now()

		cfg, err := config.Get()
		So(err, ShouldBeNil)

		var sleptFor time.Duration
		fakeSleep := func(d time.Duration) {
			sleptFor = d.Round(time.Second)
		}

		cfg.SleepFunc = fakeSleep

		clientList := clients.ClientList{
			CloudflareCacheClient: &mockClients.CloudflareCacheClienterMock{
				PurgeFunc: func(ctx context.Context, params cache.CachePurgeParams, opts ...option.RequestOption) (*cache.CachePurgeResponse, error) {
					return nil, nil
				},
			},
			LegacyCacheClient: &sdkMock.ClienterMock{
				GetCacheTimesFunc: func(ctx context.Context, auth sdk.Auth, options sdk.Options) (*models.CacheTimesList, apiError.Error) {
					return &models.CacheTimesList{
						Items: []*models.CacheTime{
							{
								ID:           "1",
								Path:         "/test-path1",
								CollectionID: "collection1",
							},
						},
						TotalCount: 1,
						Count:      1,
					}, nil
				},
			},
			SlackClient: &mockClients.SlackClienterMock{
				PostMessageFunc: func(channel string, options ...slack.MsgOption) (string, string, error) {
					return "", "", nil
				},
			},
		}
		purgeRunner, err := NewPurgeRunner(cfg, clientList)
		So(err, ShouldBeNil)

		Convey("When Run is called", func() {
			result, err := purgeRunner.Run(context.Background(), publishTime)

			Convey("Then it should complete without error", func() {
				So(err, ShouldBeNil)
				So(result.Success, ShouldBeTrue)
				So(result.Results[0].Success, ShouldBeTrue)
				So(result.Results[0].Purges, ShouldEqual, 3)
				So(sleptFor, ShouldEqual, cfg.CachePurgeDiffTime)

				Convey("And the Cloudflare cache purge should have been called", func() {
					cloudflareClient := clientList.CloudflareCacheClient.(*mockClients.CloudflareCacheClienterMock)
					So(cloudflareClient.PurgeCalls(), ShouldHaveLength, 1)

					Convey("And the Cloudflare cache purge should have been called", func() {
						cloudflareClient := clientList.CloudflareCacheClient.(*mockClients.CloudflareCacheClienterMock)
						So(cloudflareClient.PurgeCalls(), ShouldHaveLength, 1)
					})
				})
			})
		})
	})

	Convey("Given a Runner with a page to purge the cache for, mock clients that don't error and a publish time in the future", t, func() {
		timeDiffInFuture := 2 * time.Minute
		publishTime := time.Now().Add(timeDiffInFuture)

		cfg, err := config.Get()
		So(err, ShouldBeNil)

		var sleptFor time.Duration
		fakeSleep := func(d time.Duration) {
			sleptFor = d.Round(time.Second)
		}

		cfg.SleepFunc = fakeSleep

		clientList := clients.ClientList{
			CloudflareCacheClient: &mockClients.CloudflareCacheClienterMock{
				PurgeFunc: func(ctx context.Context, params cache.CachePurgeParams, opts ...option.RequestOption) (*cache.CachePurgeResponse, error) {
					return nil, nil
				},
			},
			LegacyCacheClient: &sdkMock.ClienterMock{
				GetCacheTimesFunc: func(ctx context.Context, auth sdk.Auth, options sdk.Options) (*models.CacheTimesList, apiError.Error) {
					return &models.CacheTimesList{
						Items: []*models.CacheTime{
							{
								ID:           "1",
								Path:         "/test-path1",
								CollectionID: "collection1",
							},
						},
						TotalCount: 1,
						Count:      1,
					}, nil
				},
			},
			SlackClient: &mockClients.SlackClienterMock{
				PostMessageFunc: func(channel string, options ...slack.MsgOption) (string, string, error) {
					return "", "", nil
				},
			},
		}
		purgeRunner, err := NewPurgeRunner(cfg, clientList)
		So(err, ShouldBeNil)

		Convey("When Run is called", func() {
			result, err := purgeRunner.Run(context.Background(), publishTime)

			Convey("Then it should complete without error", func() {
				So(err, ShouldBeNil)
				So(result.Success, ShouldBeTrue)
				So(result.Results[0].Success, ShouldBeTrue)
				So(result.Results[0].Purges, ShouldEqual, 3)

				Convey("And it should have slept for the correct duration", func() {
					So(sleptFor, ShouldEqual, cfg.CachePurgeDiffTime+timeDiffInFuture)
				})
			})
		})
	})

	Convey("Given a Runner with no pages to purge the cache for", t, func() {
		var sleptFor time.Duration

		cfg := &config.Configuration{
			SleepFunc: func(d time.Duration) {
				sleptFor = d.Round(time.Second)
			},
		}

		clientList := clients.ClientList{
			CloudflareCacheClient: &mockClients.CloudflareCacheClienterMock{},
			LegacyCacheClient: &sdkMock.ClienterMock{
				GetCacheTimesFunc: func(ctx context.Context, auth sdk.Auth, options sdk.Options) (*models.CacheTimesList, apiError.Error) {
					return &models.CacheTimesList{
						Items: []*models.CacheTime{},
					}, nil
				},
			},
			SlackClient: &mockClients.SlackClienterMock{},
		}
		purgeRunner, err := NewPurgeRunner(cfg, clientList)
		So(err, ShouldBeNil)

		Convey("When Run is called", func() {
			result, err := purgeRunner.Run(context.Background(), time.Now())

			Convey("Then it should complete without error and no purges", func() {
				So(err, ShouldBeNil)
				So(result.Success, ShouldBeTrue)
				So(result.Results, ShouldBeEmpty)

				Convey("And it should not have slept", func() {
					So(sleptFor, ShouldEqual, 0)

					Convey("And the Cloudflare cache purge should not have been called", func() {
						cloudflareClient := clientList.CloudflareCacheClient.(*mockClients.CloudflareCacheClienterMock)
						So(cloudflareClient.PurgeCalls(), ShouldHaveLength, 0)
					})
				})
			})
		})
	})
}

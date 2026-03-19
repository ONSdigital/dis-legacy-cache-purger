package runner

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ONSdigital/dis-legacy-cache-purger/clients"
	mockClients "github.com/ONSdigital/dis-legacy-cache-purger/clients/mock"
	"github.com/ONSdigital/dis-legacy-cache-purger/config"
	"github.com/ONSdigital/dp-legacy-cache-api/models"
	"github.com/ONSdigital/dp-legacy-cache-api/sdk"
	apiError "github.com/ONSdigital/dp-legacy-cache-api/sdk/errors"
	sdkMock "github.com/ONSdigital/dp-legacy-cache-api/sdk/mocks"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRunnerGetCacheTimes(t *testing.T) {
	Convey("Given a Purger with a Legacy Cache API that has a cache time result", t, func() {
		cfg := &config.Configuration{}
		clientList := clients.ClientList{
			CloudflareCacheClient: &mockClients.CloudflareCacheClienterMock{},
			LegacyCacheClient: &sdkMock.ClienterMock{
				GetCacheTimesFunc: func(ctx context.Context, auth sdk.Auth, options sdk.Options) (*models.CacheTimesList, apiError.Error) {
					return &models.CacheTimesList{
						Items: []*models.CacheTime{
							{
								ID:           "1",
								CollectionID: "col-1",
							},
						},
						TotalCount: 1,
						Count:      1,
					}, nil
				},
			},
			SlackClient: &mockClients.SlackClienterMock{},
		}
		purgeRunner := &PurgeRunner{
			config:     cfg,
			clientList: clientList,
		}

		Convey("When GetCacheTimes is called", func() {
			releaseTime := time.Now().Add(1 * time.Minute)
			cacheTimes, err := purgeRunner.GetCacheTimes(context.Background(), releaseTime)

			Convey("Then it should return the expected cache times without error", func() {
				So(err, ShouldBeNil)
				So(cacheTimes, ShouldHaveLength, 1)
				So(cacheTimes[0].ID, ShouldEqual, "1")
				So(cacheTimes[0].CollectionID, ShouldEqual, "col-1")
			})
		})
	})

	Convey("Given a Purger with a Legacy Cache API that has many results", t, func() {
		cfg := &config.Configuration{}
		numCacheTimes := 300
		fullCacheTimes := generateCacheTimes(numCacheTimes)

		legacyCacheClient := &sdkMock.ClienterMock{
			GetCacheTimesFunc: func(ctx context.Context, auth sdk.Auth, options sdk.Options) (*models.CacheTimesList, apiError.Error) {
				return &models.CacheTimesList{
					Items:      fullCacheTimes[options.Offset : options.Offset+options.Limit],
					TotalCount: numCacheTimes,
					Count:      options.Limit,
				}, nil
			},
		}

		clientList := clients.ClientList{
			CloudflareCacheClient: &mockClients.CloudflareCacheClienterMock{},
			LegacyCacheClient:     legacyCacheClient,
			SlackClient:           &mockClients.SlackClienterMock{},
		}
		purgeRunner := &PurgeRunner{
			config:     cfg,
			clientList: clientList,
		}

		Convey("When GetCacheTimes is called", func() {
			releaseTime := time.Now().Add(1 * time.Minute)
			cacheTimes, err := purgeRunner.GetCacheTimes(context.Background(), releaseTime)

			Convey("Then it should paginate and fetch all cache times without error", func() {
				So(err, ShouldBeNil)
				So(len(legacyCacheClient.GetCacheTimesCalls()), ShouldEqual, numCacheTimes/100)
				So(cacheTimes, ShouldHaveLength, numCacheTimes)
				So(cacheTimes[0], ShouldEqual, fullCacheTimes[0])
				So(cacheTimes[numCacheTimes-1], ShouldEqual, fullCacheTimes[numCacheTimes-1])
			})
		})
	})

	Convey("Given a Purger with a Legacy Cache API that has no cache times", t, func() {
		cfg := &config.Configuration{}
		clientList := clients.ClientList{
			CloudflareCacheClient: &mockClients.CloudflareCacheClienterMock{},
			LegacyCacheClient: &sdkMock.ClienterMock{
				GetCacheTimesFunc: func(ctx context.Context, auth sdk.Auth, options sdk.Options) (*models.CacheTimesList, apiError.Error) {
					return &models.CacheTimesList{
						Items:      []*models.CacheTime{},
						TotalCount: 0,
						Count:      0,
					}, nil
				},
			},
			SlackClient: &mockClients.SlackClienterMock{},
		}
		purgeRunner := &PurgeRunner{
			config:     cfg,
			clientList: clientList,
		}

		Convey("When GetCacheTimes is called", func() {
			releaseTime := time.Now().Add(1 * time.Minute)
			cacheTimes, err := purgeRunner.GetCacheTimes(context.Background(), releaseTime)

			Convey("Then it should return no cache times without error", func() {
				So(err, ShouldBeNil)
				So(cacheTimes, ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given a Purger with a Legacy Cache API that returns an error", t, func() {
		cfg := &config.Configuration{}
		clientList := clients.ClientList{
			CloudflareCacheClient: &mockClients.CloudflareCacheClienterMock{},
			LegacyCacheClient: &sdkMock.ClienterMock{
				GetCacheTimesFunc: func(ctx context.Context, auth sdk.Auth, options sdk.Options) (*models.CacheTimesList, apiError.Error) {
					return nil, apiError.StatusError{
						Err: fmt.Errorf("unknown error occurred"),
					}
				},
			},
			SlackClient: &mockClients.SlackClienterMock{},
		}
		purgeRunner := &PurgeRunner{
			config:     cfg,
			clientList: clientList,
		}

		Convey("When GetCacheTimes is called", func() {
			releaseTime := time.Now().Add(1 * time.Minute)
			cacheTimes, err := purgeRunner.GetCacheTimes(context.Background(), releaseTime)

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(cacheTimes, ShouldHaveLength, 0)
			})
		})
	})
}

func generateCacheTimes(num int) []*models.CacheTime {
	var cacheTimes []*models.CacheTime
	for i := 0; i < num; i++ {
		cacheTime := &models.CacheTime{
			ID: fmt.Sprintf("cache-time-%d", i),
		}
		cacheTimes = append(cacheTimes, cacheTime)
	}
	return cacheTimes
}

package runner

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ONSdigital/dis-legacy-cache-purger/clients"
	mockClients "github.com/ONSdigital/dis-legacy-cache-purger/clients/mock"
	"github.com/ONSdigital/dis-legacy-cache-purger/config"
	"github.com/cloudflare/cloudflare-go/v6/cache"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/slack-go/slack"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRunnerCollectionCachePurge(t *testing.T) {
	Convey("Given a purgeRunner with a mock Cloudflare client that does not error", t, func() {
		mockCloudflareClient := &mockClients.CloudflareCacheClienterMock{
			PurgeFunc: func(ctx context.Context, params cache.CachePurgeParams, opts ...option.RequestOption) (*cache.CachePurgeResponse, error) {
				return &cache.CachePurgeResponse{}, nil
			},
		}
		mockSlackClient := &mockClients.SlackClienterMock{
			PostMessageFunc: func(channel string, options ...slack.MsgOption) (string, string, error) {
				return "", "", nil
			},
		}

		purger := &PurgeRunner{
			clientList: clients.ClientList{
				CloudflareCacheClient: mockCloudflareClient,
				SlackClient:           mockSlackClient,
			},
			config: &config.Configuration{
				MaxParallel:         1,
				CloudflareBatchSize: 100,
			},
		}

		Convey("When CachePurgeCollection is called with valid prefixes and files", func() {
			req := CollectionCachePurgeRequest{
				CollectionID: "col-1",
				Prefixes:     []string{"/path/prefix1/", "/path/prefix2/"},
				Files:        []string{"/path/file1.html", "/path/file2.html"},
			}
			releaseTime := time.Now()
			result := purger.CachePurgeCollection(context.Background(), req, releaseTime)

			Convey("Then it should return a successful CachePurgeResult", func() {
				So(result.Success, ShouldBeTrue)
				So(result.Purges, ShouldEqual, 4)
				So(result.Error, ShouldBeNil)

				Convey("And the Cloudflare client's Purge method should have been called 2 times", func() {
					So(mockCloudflareClient.PurgeCalls(), ShouldHaveLength, 2)

					Convey("And the Slack client should have been called 1 time", func() {
						So(mockSlackClient.PostMessageCalls(), ShouldHaveLength, 1)
					})
				})
			})
		})

		Convey("When CachePurgeCollection is called with more than 100 prefixes", func() {
			numPrefixes := 101

			req := CollectionCachePurgeRequest{
				CollectionID: "col-1",
				Prefixes:     generatePrefixes(numPrefixes),
			}
			releaseTime := time.Now()
			result := purger.CachePurgeCollection(context.Background(), req, releaseTime)

			Convey("Then it should return a successful CachePurgeResult", func() {
				So(result.Success, ShouldBeTrue)
				So(result.Purges, ShouldEqual, numPrefixes)
				So(result.Error, ShouldBeNil)

				Convey("And the Cloudflare client's Purge method should have been called once for every 100 prefixes", func() {
					So(mockCloudflareClient.PurgeCalls(), ShouldHaveLength, (numPrefixes+99)/100)

					Convey("And the Slack client should have been called 1 time", func() {
						So(mockSlackClient.PostMessageCalls(), ShouldHaveLength, 1)
					})
				})
			})
		})

		Convey("When CachePurgeCollection is called with more than 100 prefixes and files", func() {
			numPurges := 101

			req := CollectionCachePurgeRequest{
				CollectionID: "col-1",
				Prefixes:     generatePrefixes(numPurges),
				Files:        generateFiles(numPurges),
			}
			releaseTime := time.Now()
			result := purger.CachePurgeCollection(context.Background(), req, releaseTime)

			Convey("Then it should return a successful CachePurgeResult", func() {
				So(result.Success, ShouldBeTrue)
				So(result.Purges, ShouldEqual, numPurges*2)
				So(result.Error, ShouldBeNil)

				Convey("And the Cloudflare client's Purge method should have been called once for every 100 prefixes and every 100 files", func() {
					So(mockCloudflareClient.PurgeCalls(), ShouldHaveLength, ((numPurges+99)/100)*2)

					Convey("And the Slack client should have been called 1 time", func() {
						So(mockSlackClient.PostMessageCalls(), ShouldHaveLength, 1)
					})
				})
			})
		})
	})

	Convey("Given a purgeRunner with a mock Cloudflare client that does error", t, func() {
		mockCloudflareClient := &mockClients.CloudflareCacheClienterMock{
			PurgeFunc: func(ctx context.Context, params cache.CachePurgeParams, opts ...option.RequestOption) (*cache.CachePurgeResponse, error) {
				return nil, fmt.Errorf("mock purge error")
			},
		}
		mockSlackClient := &mockClients.SlackClienterMock{
			PostMessageFunc: func(channel string, options ...slack.MsgOption) (string, string, error) {
				return "", "", nil
			},
		}

		purger := &PurgeRunner{
			clientList: clients.ClientList{
				CloudflareCacheClient: mockCloudflareClient,
				SlackClient:           mockSlackClient,
			},
			config: &config.Configuration{
				MaxParallel:         1,
				CloudflareBatchSize: 100,
			},
		}

		Convey("When CachePurgeCollection is called with valid prefixes and files", func() {
			req := CollectionCachePurgeRequest{
				CollectionID: "col-1",
				Prefixes:     []string{"/path/prefix1/", "/path/prefix2/"},
				Files:        []string{"/path/file1.html", "/path/file2.html"},
			}
			releaseTime := time.Now()
			result := purger.CachePurgeCollection(context.Background(), req, releaseTime)

			Convey("Then it should return a failed CachePurgeResult", func() {
				So(result.Success, ShouldBeFalse)
				So(result.Purges, ShouldEqual, 0)
				So(result.Error, ShouldNotBeNil)

				Convey("And the Cloudflare client's Purge method should have been called 2 times", func() {
					So(mockCloudflareClient.PurgeCalls(), ShouldHaveLength, 2)

					Convey("And the Slack client should have been called 1 time", func() {
						So(mockSlackClient.PostMessageCalls(), ShouldHaveLength, 1)
					})
				})
			})
		})
	})
}

func TestRunnerCachePurge(t *testing.T) {
	Convey("Given a purgeRunner with a mock Cloudflare client that does not error", t, func() {
		mockCloudflareClient := &mockClients.CloudflareCacheClienterMock{
			PurgeFunc: func(ctx context.Context, params cache.CachePurgeParams, opts ...option.RequestOption) (*cache.CachePurgeResponse, error) {
				return &cache.CachePurgeResponse{}, nil
			},
		}
		mockSlackClient := &mockClients.SlackClienterMock{
			PostMessageFunc: func(channel string, options ...slack.MsgOption) (string, string, error) {
				return "", "", nil
			},
		}

		purger := &PurgeRunner{
			clientList: clients.ClientList{
				CloudflareCacheClient: mockCloudflareClient,
				SlackClient:           mockSlackClient,
			},
			config: &config.Configuration{
				MaxParallel:         1,
				CloudflareBatchSize: 100,
			},
		}

		Convey("When CachePurge is called with multiple collections", func() {
			reqs := []CollectionCachePurgeRequest{
				{
					CollectionID: "col-1",
					Prefixes:     []string{"/path/prefix1/", "/path/prefix2/"},
					Files:        []string{"/path/file1.html", "/path/file2.html"},
				},
				{
					CollectionID: "col-2",
					Prefixes:     []string{"/path/prefix3/"},
					Files:        []string{"/path/file3.html"},
				},
			}
			releaseTime := time.Now()
			results := purger.CachePurge(context.Background(), reqs, releaseTime)

			Convey("Then it should return a successful PurgeRunResult", func() {
				So(results, ShouldHaveLength, 2)

				Convey("And each CachePurgeResult should be successful with the expected number of purges", func() {
					So(results[0].Success, ShouldBeTrue)
					So(results[0].Purges, ShouldEqual, 4)
					So(results[0].Error, ShouldBeNil)

					So(results[1].Success, ShouldBeTrue)
					So(results[1].Purges, ShouldEqual, 2)
					So(results[1].Error, ShouldBeNil)

					Convey("And the Cloudflare client's Purge method should have been called for each collection for each type of purge", func() {
						So(mockCloudflareClient.PurgeCalls(), ShouldHaveLength, 4)

						Convey("And the Slack client should be notified for each collection", func() {
							So(len(mockSlackClient.PostMessageCalls()), ShouldEqual, 2)
						})
					})
				})
			})
		})
	})
}

func TestRunnerCachePurgeByType(t *testing.T) {
	Convey("Given a purgeRunner with a mock Cloudflare client that does not error", t, func() {
		mockCloudflareClient := &mockClients.CloudflareCacheClienterMock{
			PurgeFunc: func(ctx context.Context, params cache.CachePurgeParams, opts ...option.RequestOption) (*cache.CachePurgeResponse, error) {
				return &cache.CachePurgeResponse{}, nil
			},
		}

		purger := &PurgeRunner{
			clientList: clients.ClientList{
				CloudflareCacheClient: mockCloudflareClient,
			},
			config: &config.Configuration{
				MaxParallel:         1,
				CloudflareBatchSize: 100,
			},
		}

		Convey("When CachePurgePrefixes is called with multiple prefixes", func() {
			prefixes := []string{"/path/1", "/path/2", "/path/3"}
			err := purger.CachePurgePrefixes(context.Background(), prefixes)
			Convey("Then it should complete without error", func() {
				So(err, ShouldBeNil)

				Convey("And the Cloudflare client's Purge method should have been called once with a prefix purge", func() {
					So(mockCloudflareClient.PurgeCalls(), ShouldHaveLength, 1)
					So(mockCloudflareClient.PurgeCalls()[0].Params.Body, ShouldHaveSameTypeAs, cache.CachePurgeParamsBodyCachePurgeFlexPurgeByPrefixes{})
					So(mockCloudflareClient.PurgeCalls()[0].Params.Body.(cache.CachePurgeParamsBodyCachePurgeFlexPurgeByPrefixes).Prefixes.Value, ShouldEqual, prefixes)
				})
			})
		})

		Convey("When CachePurgeFiles is called with multiple files", func() {
			files := []string{"/file?uri=test1", "/file?uri=test2", "/file?uri=test3"}
			err := purger.CachePurgeFiles(context.Background(), files)
			Convey("Then it should complete without error", func() {
				So(err, ShouldBeNil)

				Convey("And the Cloudflare client's Purge method should have been called once with a file purge", func() {
					So(mockCloudflareClient.PurgeCalls(), ShouldHaveLength, 1)
					So(mockCloudflareClient.PurgeCalls()[0].Params.Body, ShouldHaveSameTypeAs, cache.CachePurgeParamsBodyCachePurgeSingleFile{})
					So(mockCloudflareClient.PurgeCalls()[0].Params.Body.(cache.CachePurgeParamsBodyCachePurgeSingleFile).Files.Value, ShouldEqual, files)
				})
			})
		})
	})
}

func generatePrefixes(num int) []string {
	prefixes := make([]string, num)
	for i := 0; i < num; i++ {
		prefixes[i] = fmt.Sprintf("/path/prefix%d/", i)
	}
	return prefixes
}

func generateFiles(num int) []string {
	files := make([]string, num)
	for i := 0; i < num; i++ {
		files[i] = fmt.Sprintf("/path/file?uri=%d", i)
	}
	return files
}

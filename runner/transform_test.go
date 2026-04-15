package runner

import (
	"context"
	"testing"

	"github.com/ONSdigital/dp-legacy-cache-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMapCacheTimeByCollectionID(t *testing.T) {
	Convey("Given a list of CacheTime objects and domains", t, func() {
		cacheTimes := []*models.CacheTime{
			{CollectionID: "col-1", Path: "/path1"},
			{CollectionID: "col-1", Path: "/path2"},
			{CollectionID: "col-2", Path: "/path3"},
		}
		domains := []string{"domain1.com", "domain2.com"}
		ctx := context.Background()

		Convey("When mapCacheTimeByCollectionID is called", func() {
			result := mapCacheTimeByCollectionID(ctx, cacheTimes, domains)

			Convey("Then it should return the expected mapping", func() {
				So(result, ShouldResemble, map[string][]string{
					"col-1": {
						"domain1.com/path1",
						"domain2.com/path1",
						"domain1.com/path2",
						"domain2.com/path2",
					},
					"col-2": {
						"domain1.com/path3",
						"domain2.com/path3",
					},
				})
			})
		})
	})

	Convey("Given an empty list of CacheTime objects", t, func() {
		cacheTimes := []*models.CacheTime{}
		domains := []string{"domain1.com", "domain2.com"}
		ctx := context.Background()

		Convey("When mapCacheTimeByCollectionID is called", func() {
			result := mapCacheTimeByCollectionID(ctx, cacheTimes, domains)

			Convey("Then it should return an empty mapping", func() {
				So(result, ShouldBeEmpty)
			})
		})
	})

	Convey("Give a list of CacheTime objects but no domains", t, func() {
		cacheTimes := []*models.CacheTime{
			{CollectionID: "col-1", Path: "/path1"},
		}
		domains := []string{}
		ctx := context.Background()

		Convey("When mapCacheTimeByCollectionID is called", func() {
			result := mapCacheTimeByCollectionID(ctx, cacheTimes, domains)

			Convey("Then it should return an empty map", func() {
				So(result, ShouldResemble, map[string][]string{})
			})
		})
	})
}

func TestMapCollectionCacheTimeMapToRequests(t *testing.T) {
	Convey("Given a collection cache time map", t, func() {
		cacheTimeMap := map[string][]string{
			"col-1": {
				"/prefix1/path1",
				"/prefix1/path2?query=1",
				"/prefix2/path3",
			},
			"col-2": {
				"/prefix3/path4?query=2",
			},
		}
		ctx := context.Background()

		Convey("When mapCollectionCacheTimeMapToRequests is called", func() {
			requests := mapCollectionCacheTimeMapToRequests(ctx, cacheTimeMap)

			Convey("Then it should return the expected CollectionCachePurgeRequests", func() {
				expected := []CollectionCachePurgeRequest{
					{
						CollectionID: "col-1",
						Files: []string{
							"https:///prefix1/path1",
							"https:///prefix1/path1/data",
							"https:///prefix1/path1/pdf",
							"https:///prefix1/path2?query=1",
							"https:///prefix2/path3",
							"https:///prefix2/path3/data",
							"https:///prefix2/path3/pdf",
						},
					},
					{
						CollectionID: "col-2",
						Prefixes:     nil,
						Files:        []string{"https:///prefix3/path4?query=2"},
					},
				}
				So(requests, ShouldResemble, expected)
			})
		})
	})
}

func TestTransformCacheTimesToCollectionCachePurgeRequests(t *testing.T) {
	Convey("Given a list of CacheTime objects and domains", t, func() {
		cacheTimes := []*models.CacheTime{
			{CollectionID: "col-1", Path: "/path1"},
			{CollectionID: "col-1", Path: "/path2?query=1"},
			{CollectionID: "col-2", Path: "/path3"},
		}
		domains := []string{"domain1.com", "domain2.com"}
		ctx := context.Background()

		Convey("When transformCacheTimesToCollectionCachePurgeRequests is called", func() {
			requests := transformCacheTimesToCollectionCachePurgeRequests(ctx, cacheTimes, domains)

			Convey("Then it should return the expected CollectionCachePurgeRequests", func() {
				expected := []CollectionCachePurgeRequest{
					{
						CollectionID: "col-1",
						Prefixes:     nil,
						Files: []string{
							"https://domain1.com/path1",
							"https://domain1.com/path1/data",
							"https://domain1.com/path1/pdf",
							"https://domain2.com/path1",
							"https://domain2.com/path1/data",
							"https://domain2.com/path1/pdf",
							"https://domain1.com/path2?query=1",
							"https://domain2.com/path2?query=1",
						},
					},
					{
						CollectionID: "col-2",
						Prefixes:     nil,
						Files: []string{
							"https://domain1.com/path3",
							"https://domain1.com/path3/data",
							"https://domain1.com/path3/pdf",
							"https://domain2.com/path3",
							"https://domain2.com/path3/data",
							"https://domain2.com/path3/pdf",
						},
					},
				}
				So(requests, ShouldContain, expected[0])
				So(requests, ShouldContain, expected[1])
			})
		})
	})
}

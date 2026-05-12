package runner

import (
	"context"
	"strconv"
	"testing"

	"github.com/ONSdigital/dp-legacy-cache-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMapCacheTimeByCollectionID(t *testing.T) {
	Convey("Given a list of CacheTime objects and domains", t, func() {
		cacheTimes := []*models.CacheTime{
			{CollectionID: generateTestCollectionID(1), Path: generateTestPath(1)},
			{CollectionID: generateTestCollectionID(1), Path: generateTestPath(2)},
			{CollectionID: generateTestCollectionID(2), Path: generateTestPath(3)},
		}
		domains := []string{generateTestDomain(1), generateTestDomain(2)}
		ctx := context.Background()

		Convey("When mapCacheTimeByCollectionID is called", func() {
			result := mapCacheTimeByCollectionID(ctx, cacheTimes, domains)

			Convey("Then it should return the expected mapping", func() {
				So(result, ShouldResemble, map[string][]string{
					generateTestCollectionID(1): {
						generateTestDomain(1) + generateTestPath(1),
						generateTestDomain(2) + generateTestPath(1),
						generateTestDomain(1) + generateTestPath(2),
						generateTestDomain(2) + generateTestPath(2),
					},
					generateTestCollectionID(2): {
						generateTestDomain(1) + generateTestPath(3),
						generateTestDomain(2) + generateTestPath(3),
					},
				})
			})
		})
	})

	Convey("Given an empty list of CacheTime objects", t, func() {
		cacheTimes := []*models.CacheTime{}
		domains := []string{generateTestDomain(1), generateTestDomain(2)}
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
			{CollectionID: generateTestCollectionID(1), Path: generateTestPath(1)},
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
			generateTestCollectionID(1): {
				"/prefix1/path1",
				"/prefix1/path2?query=1",
				"/prefix2/path3",
			},
			generateTestCollectionID(2): {
				"/prefix3/path4?query=2",
			},
		}
		ctx := context.Background()

		Convey("When mapCollectionCacheTimeMapToRequests is called", func() {
			requests := mapCollectionCacheTimeMapToRequests(ctx, cacheTimeMap)

			Convey("Then it should return the expected CollectionCachePurgeRequests", func() {
				expected := []CollectionCachePurgeRequest{
					{
						CollectionID: generateTestCollectionID(1),
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
						CollectionID: generateTestCollectionID(2),
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
			{CollectionID: generateTestCollectionID(1), Path: generateTestPath(1)},
			{CollectionID: generateTestCollectionID(1), Path: generateTestPath(2) + "?query=1"},
			{CollectionID: generateTestCollectionID(2), Path: generateTestPath(3)},
		}
		domains := []string{generateTestDomain(1), generateTestDomain(2)}
		ctx := context.Background()

		Convey("When transformCacheTimesToCollectionCachePurgeRequests is called", func() {
			requests := transformCacheTimesToCollectionCachePurgeRequests(ctx, cacheTimes, domains)

			Convey("Then it should return the expected CollectionCachePurgeRequests", func() {
				expected := []CollectionCachePurgeRequest{
					{
						CollectionID: generateTestCollectionID(1),
						Prefixes:     nil,
						Files: []string{
							"https://" + generateTestDomain(1) + generateTestPath(1),
							"https://" + generateTestDomain(1) + generateTestPath(1) + "/data",
							"https://" + generateTestDomain(1) + generateTestPath(1) + "/pdf",
							"https://" + generateTestDomain(2) + generateTestPath(1),
							"https://" + generateTestDomain(2) + generateTestPath(1) + "/data",
							"https://" + generateTestDomain(2) + generateTestPath(1) + "/pdf",
							"https://" + generateTestDomain(1) + generateTestPath(2) + "?query=1",
							"https://" + generateTestDomain(2) + generateTestPath(2) + "?query=1",
						},
					},
					{
						CollectionID: generateTestCollectionID(2),
						Prefixes:     nil,
						Files: []string{
							"https://" + generateTestDomain(1) + generateTestPath(3),
							"https://" + generateTestDomain(1) + generateTestPath(3) + "/data",
							"https://" + generateTestDomain(1) + generateTestPath(3) + "/pdf",
							"https://" + generateTestDomain(2) + generateTestPath(3),
							"https://" + generateTestDomain(2) + generateTestPath(3) + "/data",
							"https://" + generateTestDomain(2) + generateTestPath(3) + "/pdf",
						},
					},
				}
				So(requests, ShouldContain, expected[0])
				So(requests, ShouldContain, expected[1])
			})
		})
	})
}

func generateTestDomain(i int) string {
	return "domain" + strconv.Itoa(i) + ".com"
}

func generateTestPath(i int) string {
	return "/path" + strconv.Itoa(i)
}

func generateTestCollectionID(i int) string {
	return "col-" + strconv.Itoa(i)
}

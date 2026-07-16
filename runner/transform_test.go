package runner

import (
	"context"
	"strconv"
	"testing"

	"github.com/ONSdigital/dp-legacy-cache-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	testPath         = "/path"
	testDatasetsPath = "/datasets/"
)

func TestMapCacheTimeByCollection(t *testing.T) {
	Convey("Given a list of CacheTime objects and domains", t, func() {
		cacheTimes := []*models.CacheTime{
			{CollectionID: generateTestCollectionID(1), CollectionTitle: generateTestCollectionTitle(1), Path: generateTestPath(1)},
			{CollectionID: generateTestCollectionID(1), CollectionTitle: generateTestCollectionTitle(1), Path: generateTestPath(2)},
			{CollectionID: generateTestCollectionID(2), CollectionTitle: generateTestCollectionTitle(2), Path: generateTestPath(3)},
		}
		domains := []string{generateTestDomain(1), generateTestDomain(2)}
		ctx := context.Background()

		Convey("When mapCacheTimeByCollection is called", func() {
			result := mapCacheTimeByCollection(ctx, cacheTimes, domains)

			Convey("Then it should return the expected mapping", func() {
				So(result, ShouldResemble, map[string]CollectionCachePaths{
					generateTestCollectionID(1): {
						collectionID:    generateTestCollectionID(1),
						collectionTitle: generateTestCollectionTitle(1),
						paths: []string{
							generateTestDomain(1) + generateTestPath(1),
							generateTestDomain(2) + generateTestPath(1),
							generateTestDomain(1) + generateTestPath(2),
							generateTestDomain(2) + generateTestPath(2),
						},
					},
					generateTestCollectionID(2): {
						collectionID:    generateTestCollectionID(2),
						collectionTitle: generateTestCollectionTitle(2),
						paths: []string{
							generateTestDomain(1) + generateTestPath(3),
							generateTestDomain(2) + generateTestPath(3),
						},
					},
				})
			})
		})
	})

	Convey("Given an empty list of CacheTime objects", t, func() {
		cacheTimes := []*models.CacheTime{}
		domains := []string{generateTestDomain(1), generateTestDomain(2)}
		ctx := context.Background()

		Convey("When mapCacheTimeByCollection is called", func() {
			result := mapCacheTimeByCollection(ctx, cacheTimes, domains)

			Convey("Then it should return an empty mapping", func() {
				So(result, ShouldBeEmpty)
			})
		})
	})

	Convey("Give a list of CacheTime objects but no domains", t, func() {
		cacheTimes := []*models.CacheTime{
			{CollectionID: generateTestCollectionID(1), CollectionTitle: generateTestCollectionTitle(1), Path: generateTestPath(1)},
		}
		domains := []string{}
		ctx := context.Background()

		Convey("When mapCacheTimeByCollection is called", func() {
			result := mapCacheTimeByCollection(ctx, cacheTimes, domains)

			Convey("Then it should return an empty map", func() {
				So(result, ShouldResemble, map[string]CollectionCachePaths{})
			})
		})
	})
}

func TestMapCollectionCacheTimeMapToRequests(t *testing.T) {
	Convey("Given a collection cache time map", t, func() {
		cacheTimeMap := map[string]CollectionCachePaths{
			generateTestCollectionID(1): CollectionCachePaths{
				collectionID:    generateTestCollectionID(1),
				collectionTitle: generateTestCollectionTitle(1),
				paths: []string{
					"/prefix1/path1",
					"/prefix1/path2?query=1",
					"/prefix2/path3",
				},
			},
			generateTestCollectionID(2): CollectionCachePaths{
				collectionID:    generateTestCollectionID(2),
				collectionTitle: generateTestCollectionTitle(2),
				paths: []string{
					"/prefix3/path4?query=2",
				},
			},
		}
		ctx := context.Background()
		dataPaths := []string{"/prefix2/"}
		pdfPaths := []string{"/prefix1/"}
		excludedPaths := []string{}

		Convey("When mapCollectionCacheTimeMapToRequests is called", func() {
			requests := mapCollectionCacheTimeMapToRequests(ctx, cacheTimeMap, dataPaths, pdfPaths, excludedPaths)

			Convey("Then it should return the expected CollectionCachePurgeRequests", func() {
				expected := []CollectionCachePurgeRequest{
					{
						CollectionID:    generateTestCollectionID(1),
						CollectionTitle: generateTestCollectionTitle(1),
						Files: []string{
							"https:///prefix1/path1",
							"https:///prefix1/path1/pdf",
							"https:///prefix1/path2?query=1",
							"https:///prefix2/path3",
							"https:///prefix2/path3/data",
						},
					},
					{
						CollectionID:    generateTestCollectionID(2),
						CollectionTitle: generateTestCollectionTitle(2),
						Prefixes:        nil,
						Files:           []string{"https:///prefix3/path4?query=2"},
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
			{CollectionID: generateTestCollectionID(1), CollectionTitle: generateTestCollectionTitle(1), Path: generateTestPath(1)},
			{CollectionID: generateTestCollectionID(1), CollectionTitle: generateTestCollectionTitle(1), Path: generateTestPath(2) + "?query=1"},
			{CollectionID: generateTestCollectionID(2), CollectionTitle: generateTestCollectionTitle(2), Path: generateTestPath(3)},
		}
		domains := []string{generateTestDomain(1), generateTestDomain(2)}
		ctx := context.Background()

		Convey("When transformCacheTimesToCollectionCachePurgeRequests is called", func() {
			requests := transformCacheTimesToCollectionCachePurgeRequests(ctx, cacheTimes, domains, []string{testPath}, []string{testPath}, []string{})

			Convey("Then it should return the expected CollectionCachePurgeRequests", func() {
				expected := []CollectionCachePurgeRequest{
					{
						CollectionID:    generateTestCollectionID(1),
						CollectionTitle: generateTestCollectionTitle(1),
						Prefixes:        nil,
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
						CollectionID:    generateTestCollectionID(2),
						CollectionTitle: generateTestCollectionTitle(2),
						Prefixes:        nil,
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

func TestPathToFilesExcludedPaths(t *testing.T) {
	Convey("Given a path that matches an excluded pattern", t, func() {
		excludedPaths := []string{"/timeseries/", "/previous/"}

		Convey("When the path contains /previous/", func() {
			result := pathToFiles("domain.com/datasets/some-dataset/current/previous/v131", []string{testDatasetsPath}, []string{}, excludedPaths)

			Convey("Then it should return nil", func() {
				So(result, ShouldBeNil)
			})
		})

		Convey("When the path contains /timeseries/", func() {
			result := pathToFiles("domain.com/economy/timeseries/abc123", []string{testDatasetsPath}, []string{}, excludedPaths)

			Convey("Then it should return nil", func() {
				So(result, ShouldBeNil)
			})
		})

		Convey("When the path does not match any excluded pattern", func() {
			result := pathToFiles("domain.com/datasets/some-dataset/current", []string{testDatasetsPath}, []string{}, excludedPaths)

			Convey("Then it should return file entries", func() {
				So(result, ShouldNotBeNil)
				So(len(result), ShouldEqual, 2)
			})
		})
	})
}

func generateTestDomain(i int) string {
	return "domain" + strconv.Itoa(i) + ".com"
}

func generateTestPath(i int) string {
	return testPath + strconv.Itoa(i)
}

func generateTestCollectionID(i int) string {
	return "col-" + strconv.Itoa(i)
}

func generateTestCollectionTitle(i int) string {
	return "Collection " + strconv.Itoa(i)
}

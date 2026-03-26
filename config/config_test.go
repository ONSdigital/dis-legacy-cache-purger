package config

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	Convey("Given an environment with no environment variables set", t, func() {
		os.Clearenv()
		cfg, err := Get()

		Convey("When the config values are retrieved", func() {
			Convey("Then there should be no error returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("The values should be set to the expected defaults", func() {
				So(cfg.CachePurgeDiffTime, ShouldEqual, 30*time.Second)
				So(cfg.CloudflareAPIToken, ShouldEqual, "")
				So(cfg.CloudflareBatchSize, ShouldEqual, 100)
				So(cfg.CloudflareZoneID, ShouldEqual, "")
				So([]string(cfg.Domains), ShouldResemble, []string{"sandbox.onsdigital.co.uk"})
				So(cfg.EnableCloudflarePurge, ShouldBeFalse)
				So(cfg.EnableCacheAPI, ShouldBeFalse)
				So(cfg.EnableSlackAlerts, ShouldBeFalse)
				So(cfg.LegacyCacheAPIURL, ShouldEqual, "http://localhost:29100")
				So(cfg.MaxParallel, ShouldEqual, 10)
				So(cfg.ServiceToken, ShouldEqual, "cache-purger-test-auth-token")
				So(cfg.SlackAPIToken, ShouldEqual, "")
				So(cfg.SlackChannel, ShouldEqual, "#sandbox-publish-log")
			})
		})
	})
}

func TestSensitiveFieldsOmitted(t *testing.T) {
	Convey("Given an environment with no environment variables set", t, func() {
		os.Clearenv()
		cfg, err := Get()

		Convey("When string is called the service token is not returned", func() {
			obj := cfg.String()
			So(err, ShouldBeNil)
			So(obj, ShouldNotContainSubstring, cfg.ServiceToken)
		})
	})
}

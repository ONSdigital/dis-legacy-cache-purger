package config

import (
	"encoding/json"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dis-legacy-cache-purger
type Configuration struct {
	CachePurgeDiffTime    time.Duration       `envconfig:"CACHE_PURGE_DIFF_TIME"`
	CloudflareAPIToken    string              `envconfig:"CLOUDFLARE_API_TOKEN" json:"-"`
	CloudflareBatchSize   int                 `envconfig:"CLOUDFLARE_BATCH_SIZE"`
	CloudflareZoneID      string              `envconfig:"CLOUDFLARE_ZONE_ID" json:"-"`
	Domains               []string            `envconfig:"DOMAINS" json:"domains"`
	EnableCloudflarePurge bool                `envconfig:"ENABLE_CLOUDFLARE_PURGE" json:"enable_cloudflare_purge"`
	EnableCacheAPI        bool                `envconfig:"ENABLE_CACHE_API" json:"enable_cache_api"`
	EnableSlackAlerts     bool                `envconfig:"ENABLE_SLACK_ALERTS" json:"enable_slack_alerts"`
	LegacyCacheAPIURL     string              `envconfig:"LEGACY_CACHE_API_URL"`
	MaxParallel           int                 `envconfig:"MAX_PARALLEL" json:"max_parallel"`
	ServiceToken          string              `envconfig:"LEGACY_CACHE_API_SERVICE_TOKEN"  json:"-"`
	SlackAPIToken         string              `envconfig:"SLACK_API_TOKEN" json:"-"`
	SlackChannel          string              `envconfig:"SLACK_CHANNEL" json:"slack_channel"`
	SleepFunc             func(time.Duration) `envconfig:"-" json:"-"`
}

var cfg *Configuration

// Get returns the config with variables loaded from environment variables
func Get() (*Configuration, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Configuration{
		CachePurgeDiffTime:    30 * time.Second,
		CloudflareAPIToken:    "",
		CloudflareBatchSize:   100,
		CloudflareZoneID:      "",
		Domains:               []string{"sandbox.onsdigital.co.uk"},
		EnableCloudflarePurge: false,
		EnableCacheAPI:        false,
		EnableSlackAlerts:     false,
		ServiceToken:          "cache-purger-test-auth-token",
		SlackAPIToken:         "",
		SlackChannel:          "#sandbox-publish-log",
		LegacyCacheAPIURL:     "http://localhost:29100",
		MaxParallel:           10, // default value
		SleepFunc: func(d time.Duration) {
			time.Sleep(d)
		},
	}

	return cfg, envconfig.Process("", cfg)
}

// String is implemented to prevent sensitive fields being logged.
// The config is returned as JSON with sensitive fields omitted.
func (config Configuration) String() string {
	b, _ := json.Marshal(config)
	return string(b)
}

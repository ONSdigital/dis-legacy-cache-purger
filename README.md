# dis-legacy-cache-purger

dis-legacy-cache-purger - a scheduler application for automating the purging of cache for paths to be published.

## Getting started

* Run `make debug`

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable           | Default                        | Description                                                         |
|--------------------------------|--------------------------------|---------------------------------------------------------------------|
| CACHE_PURGE_DIFF_TIME          | 30s                            | Time to wait before purging cache after publish (difference window) |
| CLOUDFLARE_API_TOKEN           | ""                             | The API token for Cloudflare                                        |
| CLOUDFLARE_BATCH_SIZE          | 100                            | Number of paths per batch for Cloudflare purge                      |
| CLOUDFLARE_ZONE_ID             | ""                             | The Cloudflare Zone ID                                              |
| DOMAINS                        | ["sandbox.onsdigital.co.uk"]   | List of domains to use for cache purging                            |
| ENABLE_CACHE_API               | false                          | Enable use of the legacy cache API                                  |
| ENABLE_CLOUDFLARE_PURGE        | false                          | Enable Cloudflare cache purging                                     |
| ENABLE_SLACK_ALERTS            | false                          | Enable Slack alert notifications                                    |
| LEGACY_CACHE_API_SERVICE_TOKEN | "cache-purger-test-auth-token" | The service auth token to connect to dp-legacy-cache-api            |
| LEGACY_CACHE_API_URL           | "http://localhost:29100"       | The URL for dp-legacy-cache-api                                     |
| MAX_PARALLEL                   | 10                             | Maximum number of parallel operations                               |
| SLACK_API_TOKEN                | ""                             | The API token for Slack                                             |
| SLACK_CHANNEL                  | "#sandbox-publish-log"         | The Slack channel to send notifications to                          |

### Tools

To run some of our tests you will need additional tooling:

#### Audit

We use `dis-vulncheck` to do auditing, which you will [need to install](https://github.com/ONSdigital/dis-vulncheck).

#### Linting

We use v2 of golangci-lint, which you will [need to install](https://golangci-lint.run/docs/welcome/install).

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

## License

Copyright © 2026, Office for National Statistics (<https://www.ons.gov.uk>)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

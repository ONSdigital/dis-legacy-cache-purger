package clients

import (
	"context"

	"github.com/ONSdigital/dp-legacy-cache-api/sdk"
	"github.com/cloudflare/cloudflare-go/v6/cache"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/slack-go/slack"
)

//go:generate moq -out mock/cache.go -pkg mock . CloudflareCacheClienter
//go:generate moq -out mock/slack.go -pkg mock . SlackClienter

type CloudflareCacheClienter interface {
	Purge(ctx context.Context, params cache.CachePurgeParams, opts ...option.RequestOption) (*cache.CachePurgeResponse, error)
}

type SlackClienter interface {
	PostMessage(channel string, options ...slack.MsgOption) (string, string, error)
}

// clients.ClientList is a struct obj of all the clients the service is dependent on
type ClientList struct {
	CloudflareCacheClient CloudflareCacheClienter
	LegacyCacheClient     sdk.Clienter
	SlackClient           SlackClienter
}

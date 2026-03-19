package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/slack-go/slack"
)

func (p *PurgeRunner) sendSuccessMessageForCollection(ctx context.Context, req CollectionCachePurgeRequest, releaseTime time.Time) error {
	log.Info(ctx, "sending success slack notification for collection", log.Data{"collection_id": req.CollectionID})

	pathCount := len(req.Prefixes) + len(req.Files)

	fields := []*slack.TextBlockObject{
		slack.NewTextBlockObject("mrkdwn", "*Collection*", false, false),
		slack.NewTextBlockObject("mrkdwn", "*Publish date*", false, false),
		slack.NewTextBlockObject("plain_text", req.CollectionID, false, false),
		slack.NewTextBlockObject("plain_text", releaseTime.Format(time.RFC3339), false, false),
		slack.NewTextBlockObject("mrkdwn", "*Number of paths purged*", false, false),
		slack.NewTextBlockObject("mrkdwn", "*Time of cache purge*", false, false),
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("%d", pathCount), false, false),
		slack.NewTextBlockObject("mrkdwn", time.Now().UTC().Format(time.RFC3339), false, false),
	}

	section := slack.NewSectionBlock(
		nil,
		fields,
		nil,
	)

	attachment := slack.Attachment{
		Color: "#2eb886",
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{section},
		},
	}

	_, _, err := p.clientList.SlackClient.PostMessage(
		p.config.SlackChannel,
		slack.MsgOptionText("Cache purged for collection", false),
		slack.MsgOptionAttachments(attachment),
	)
	return err
}

func (p *PurgeRunner) sendFailureMessageForCollection(ctx context.Context, req CollectionCachePurgeRequest, releaseTime time.Time) error {
	log.Info(ctx, "sending failure slack notification for collection", log.Data{"collection_id": req.CollectionID})

	// TODO: abstract this to reduce duplication with success message
	pathCount := len(req.Prefixes) + len(req.Files)

	fields := []*slack.TextBlockObject{
		slack.NewTextBlockObject("mrkdwn", "*Collection*", false, false),
		slack.NewTextBlockObject("mrkdwn", "*Publish date*", false, false),
		slack.NewTextBlockObject("plain_text", req.CollectionID, false, false),
		slack.NewTextBlockObject("plain_text", releaseTime.Format(time.RFC3339), false, false),
		slack.NewTextBlockObject("mrkdwn", "*Number of paths purged*", false, false),
		slack.NewTextBlockObject("mrkdwn", "*Time of cache purge*", false, false),
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("%d", pathCount), false, false),
		slack.NewTextBlockObject("mrkdwn", time.Now().UTC().Format(time.RFC3339), false, false),
	}

	section := slack.NewSectionBlock(
		nil,

		fields,
		nil,
	)

	attachment := slack.Attachment{
		Color: "#e01e5a",
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{section},
		},
	}

	_, _, err := p.clientList.SlackClient.PostMessage(
		p.config.SlackChannel,
		slack.MsgOptionText("Cache purge failed for collection", false),
		slack.MsgOptionAttachments(attachment),
	)
	return err
}

package runner

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/slack-go/slack"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dis-legacy-cache-purger/clients"
	mockClients "github.com/ONSdigital/dis-legacy-cache-purger/clients/mock"
	"github.com/ONSdigital/dis-legacy-cache-purger/config"
)

const (
	testSlackChannel = "#test"
)

func TestSendSuccessAndFailureMessageForCollection(t *testing.T) {
	Convey("Given a PurgeRunner with a mock Slack client that does not return an error", t, func() {
		mockSlackClient := &mockClients.SlackClienterMock{
			PostMessageFunc: func(channel string, options ...slack.MsgOption) (string, string, error) {
				return "ts", channel, nil
			},
		}
		purger := &PurgeRunner{
			clientList: clients.ClientList{
				SlackClient: mockSlackClient,
			},
			config: &config.Configuration{
				SlackChannel: testSlackChannel,
			},
		}

		req := CollectionCachePurgeRequest{
			CollectionID: generateTestCollectionID(1),
			Prefixes:     []string{generateTestPath(1), generateTestPath(2)},
			Files:        []string{generateTestPath(3), generateTestPath(4)},
		}
		releaseTime := time.Now()

		Convey("When sendSuccessMessageForCollection is called", func() {
			err := purger.sendSuccessMessageForCollection(context.Background(), req, releaseTime)
			Convey("Then it should not return an error", func() {
				So(err, ShouldBeNil)
				So(len(mockSlackClient.PostMessageCalls()), ShouldEqual, 1)
				call := mockSlackClient.PostMessageCalls()[0]
				So(call.Channel, ShouldEqual, testSlackChannel)
			})
		})

		Convey("When sendFailureMessageForCollection is called", func() {
			err := purger.sendFailureMessageForCollection(context.Background(), req, releaseTime)
			Convey("Then it should not return an error", func() {
				So(err, ShouldBeNil)
				So(len(mockSlackClient.PostMessageCalls()), ShouldEqual, 1)
				call := mockSlackClient.PostMessageCalls()[0]
				So(call.Channel, ShouldEqual, testSlackChannel)
			})
		})
	})
}

func TestSendSuccessMessageForCollection_Error(t *testing.T) {
	Convey("Given a PurgeRunner with a mock Slack client that returns an error", t, func() {
		mockSlackClient := &mockClients.SlackClienterMock{
			PostMessageFunc: func(channel string, options ...slack.MsgOption) (string, string, error) {
				return "", channel, fmt.Errorf("mock error")
			},
		}
		purger := &PurgeRunner{
			clientList: clients.ClientList{
				SlackClient: mockSlackClient,
			},
			config: &config.Configuration{
				SlackChannel: testSlackChannel,
			},
		}
		req := CollectionCachePurgeRequest{
			CollectionID: generateTestCollectionID(1),
			Prefixes:     []string{generateTestPath(1)},
			Files:        []string{generateTestPath(2)},
		}
		releaseTime := time.Now()

		Convey("When sendSuccessMessageForCollection is called and Slack errors", func() {
			err := purger.sendSuccessMessageForCollection(context.Background(), req, releaseTime)
			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When sendFailureMessageForCollection is called and Slack errors", func() {
			err := purger.sendFailureMessageForCollection(context.Background(), req, releaseTime)
			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

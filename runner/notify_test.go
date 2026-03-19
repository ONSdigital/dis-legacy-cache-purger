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
				SlackChannel: "#test",
			},
		}

		req := CollectionCachePurgeRequest{
			CollectionID: "col-1",
			Prefixes:     []string{"/prefix1", "/prefix2"},
			Files:        []string{"/file1", "/file2"},
		}
		releaseTime := time.Now()

		Convey("When sendSuccessMessageForCollection is called", func() {
			err := purger.sendSuccessMessageForCollection(context.Background(), req, releaseTime)
			Convey("Then it should not return an error", func() {
				So(err, ShouldBeNil)
				So(len(mockSlackClient.PostMessageCalls()), ShouldEqual, 1)
				call := mockSlackClient.PostMessageCalls()[0]
				So(call.Channel, ShouldEqual, "#test")
			})
		})

		Convey("When sendFailureMessageForCollection is called", func() {
			err := purger.sendFailureMessageForCollection(context.Background(), req, releaseTime)
			Convey("Then it should not return an error", func() {
				So(err, ShouldBeNil)
				So(len(mockSlackClient.PostMessageCalls()), ShouldEqual, 1)
				call := mockSlackClient.PostMessageCalls()[0]
				So(call.Channel, ShouldEqual, "#test")
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
				SlackChannel: "#test",
			},
		}
		req := CollectionCachePurgeRequest{
			CollectionID: "col-err",
			Prefixes:     []string{"/prefix1"},
			Files:        []string{"/file1"},
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

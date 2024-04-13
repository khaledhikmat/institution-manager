package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/khaledhikmat/institution-manager/shared/equates"
	"github.com/khaledhikmat/institution-manager/shared/service/campaign"
	"github.com/khaledhikmat/institution-manager/shared/service/member"

	"github.com/mitchellh/mapstructure"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

var pledgesTopicSubscription = &common.Subscription{
	PubsubName: equates.CAMPAIGN_PUB_SUB,
	Topic:      equates.PLEDGES_TOPIC,
	Route:      fmt.Sprintf("/%s", equates.PLEDGES_TOPIC),
}

// Global DAPR client
var canxCtx context.Context
var daprclient dapr.Client
var cmpService campaign.IService

func main() {
	rootCtx := context.Background()
	canxCtx, _ = signal.NotifyContext(rootCtx, os.Interrupt)

	// Create a DAPR service using a hard-coded port (must match start-notifier.sh)
	s := daprd.NewService(":8081")
	fmt.Println("DAPR Service created!")

	// Create a DAPR client
	// Must be a global client since it is singleton
	// Hence it would be injected in actor packages as needed
	c, err := dapr.NewClient()
	if err != nil {
		fmt.Println("Failed to start dapr client", err)
		return
	}
	daprclient = c
	defer daprclient.Close()

	cmpService = campaign.NewService(canxCtx)
	defer cmpService.Finalize()

	// Register pub/sub pledge handlers
	if err := s.AddTopicEventHandler(pledgesTopicSubscription, pledgeNotifierHandler); err != nil {
		panic(err)
	}
	fmt.Println("Pledges topic handler registered!")

	// Start DAPR service
	// TODO: Provide cancellation context
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func pledgeNotifierHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	fmt.Println("pledgeNotifierHandler....")

	go func() {
		// Decode pledge
		evt := member.MemberPledge{}
		mapstructure.Decode(e.Data, &evt)

		fmt.Printf("Received an notifier pledge for CAMPAIGN ID %s - MEMBER ID %s - amount %d\n",
			evt.CampaignID,
			evt.MemberID,
			evt.Amount)

		// TODO: Notify
	}()

	return false, nil
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/khaledhikmat/institution-manager/shared/equates"
	"github.com/khaledhikmat/institution-manager/shared/service/campaign"
	"github.com/khaledhikmat/institution-manager/shared/service/realtimer"

	"github.com/mitchellh/mapstructure"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

var campaignsTopicSubscription = &common.Subscription{
	PubsubName: equates.CAMPAIGN_PUB_SUB,
	Topic:      equates.CAMPAIGNS_TOPIC,
	Route:      fmt.Sprintf("/%s", equates.CAMPAIGNS_TOPIC),
}

// Global DAPR client
var canxCtx context.Context
var daprclient dapr.Client
var cmpService campaign.IService
var realtimeService realtimer.IService

func main() {
	rootCtx := context.Background()
	canxCtx, _ = signal.NotifyContext(rootCtx, os.Interrupt)

	// Load env vars
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to start env vars", err)
		return
	}

	// Create a DAPR service using a hard-coded port (must match start-externalizer.sh)
	s := daprd.NewService(":8083")
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

	realtimeService, err = realtimer.NewAblyService(canxCtx, os.Getenv("ABLY_API_KEY"), "institution-manager.externalyzer")
	if err != nil {
		fmt.Println("Failed to start ably client", err)
		return
	}
	defer realtimeService.Finalize()

	cmpService = campaign.NewService(canxCtx)
	defer cmpService.Finalize()

	// Register pub/sub campaigns handlers
	if err := s.AddTopicEventHandler(campaignsTopicSubscription, campaignExternalizerHandler); err != nil {
		panic(err)
	}
	fmt.Println("Campaigns topic handler registered!")

	// Start DAPR service
	// TODO: Provide cancellation context
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func campaignExternalizerHandler(_ context.Context, e *common.TopicEvent) (retry bool, err error) {
	fmt.Println("campaignExternalizerHandler....", e.Data)

	go func() {
		// Decode pledge
		// evt := equates.CampaignPledges{}
		evt := campaign.Campaign{}
		err := mapstructure.Decode(e.Data, &evt)
		if err != nil {
			fmt.Println("Failed to decode event", err)
			return
		}

		fmt.Printf("Received an externalizer pledge for CAMPAIGN %s\n",
			evt.Name)

		// Send to Realtime service
		err = realtimeService.Externalize("campaign-"+evt.ID, "campaign-update", evt)
		if err != nil {
			fmt.Printf("error publishing %v", err)
		}
	}()

	return false, nil
}

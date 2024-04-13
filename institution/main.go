package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	institutionactor "github.com/khaledhikmat/institution-manager/shared/actors/institution"
	"github.com/khaledhikmat/institution-manager/shared/equates"
	"github.com/khaledhikmat/institution-manager/shared/service/campaign"
	"github.com/khaledhikmat/institution-manager/shared/service/institution"
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
var insService institution.IService

func main() {
	rootCtx := context.Background()
	canxCtx, _ = signal.NotifyContext(rootCtx, os.Interrupt)

	// Create a DAPR service using a hard-coded port (must match start-institution.sh)
	s := daprd.NewService(":8082")
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
	insService = institution.NewService(canxCtx)
	defer insService.Finalize()

	// Inject dapr client and other services in actor packages
	institutionactor.Daprclient = daprclient
	institutionactor.InstitutionService = insService

	// Register actors
	s.RegisterActorImplFactoryContext(institutionactor.InstitutionActorFactory)
	fmt.Println("Institution actor registered!")

	// Register pub/sub pledge handlers
	if err := s.AddTopicEventHandler(pledgesTopicSubscription, pledgeInstitutionHandler); err != nil {
		panic(err)
	}
	fmt.Println("Pledges topic handler registered!")

	// Start DAPR service
	// TODO: Provide cancellation context
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func pledgeInstitutionHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	fmt.Println("pledgeInstitutionHandler....")

	go func() {
		// Decode pledge
		evt := member.MemberPledge{}
		mapstructure.Decode(e.Data, &evt)

		fmt.Printf("Received an institution pledge for campaign %s - amount %d\n",
			evt.CampaignID,
			evt.Amount)

		cmp, err := cmpService.GetCampaign(evt.CampaignID)
		if err != nil {
			fmt.Println("pledgeInstitutionHandler get campaign error", err)
			return
		}

		ins, err := insService.GetInstitution(cmp.InstitutionID)
		if err != nil {
			fmt.Println("pledgeInstitutionHandler get institution error", err)
			return
		}

		// Resolve actor by campaign institution
		insActorProxy := institutionactor.NewInstitutionActor(ins.Id)
		daprclient.ImplActorClientStub(insActorProxy)

		// Call actor methods using the main canxCtx
		err = insActorProxy.Pledge(canxCtx, evt)
		if err != nil {
			fmt.Println("pledgeInstitutionHandler calling actor error", err)
			return
		}
	}()

	return false, nil
}

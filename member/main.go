package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	memberactor "github.com/khaledhikmat/institution-manager/shared/actors/member"
	"github.com/khaledhikmat/institution-manager/shared/equates"
	"github.com/khaledhikmat/institution-manager/shared/service/member"

	"github.com/mitchellh/mapstructure"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

var pledgesTopicSubscription = &common.Subscription{
	PubsubName: equates.InstitutionManagerPubSub,
	Topic:      equates.PledgesTopic,
	Route:      fmt.Sprintf("/%s", equates.PledgesTopic),
}

// Global DAPR client
var canxCtx context.Context
var daprclient dapr.Client
var memberService member.IService

func main() {
	rootCtx := context.Background()
	canxCtx, _ = signal.NotifyContext(rootCtx, os.Interrupt)

	// Create a DAPR service using a hard-coded port (must match start-member.sh)
	s := daprd.NewService(":8084")
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

	memberService = member.NewService(canxCtx)
	defer memberService.Finalize()

	// Inject dapr client and other services in actor packages
	memberactor.Daprclient = daprclient
	memberactor.MemberService = memberService

	// Register actors
	s.RegisterActorImplFactoryContext(memberactor.MemberActorFactory)
	fmt.Println("Member actor registered!")

	// Register pub/sub pledge handlers
	if err := s.AddTopicEventHandler(pledgesTopicSubscription, pledgeMemberHandler); err != nil {
		panic(err)
	}
	fmt.Println("Pledges topic handler registered!")

	// Start DAPR service
	// TODO: Provide cancellation context
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func pledgeMemberHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	fmt.Println("pledgeMemberHandler....")

	go func() {
		// Decode pledge
		evt := member.MemberPledge{}
		mapstructure.Decode(e.Data, &evt)

		fmt.Printf("Received an member pledge for CAMPAIGN ID %s - MEMBER ID %s - amount %d\n",
			evt.CampaignID,
			evt.MemberID,
			evt.Amount)

		// Resolve actor by member id
		memActorProxy := memberactor.NewMemberActor(evt.MemberID)
		daprclient.ImplActorClientStub(memActorProxy)

		// Call actor methods using the main canxCtx
		err = memActorProxy.Pledge(canxCtx, evt)
		if err != nil {
			fmt.Println("pledgeMemberHandler calling actor error", err)
			return
		}
	}()

	return false, nil
}

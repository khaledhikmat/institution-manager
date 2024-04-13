package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	membershipactor "github.com/khaledhikmat/institution-manager/shared/actors/membership"
	"github.com/khaledhikmat/institution-manager/shared/equates"
	"github.com/khaledhikmat/institution-manager/shared/service/member"
	"github.com/khaledhikmat/institution-manager/shared/service/membership"

	"github.com/mitchellh/mapstructure"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

var transactionsTopicSubscription = &common.Subscription{
	PubsubName: equates.CAMPAIGN_PUB_SUB,
	Topic:      equates.TRANSACTIONS_TOPIC,
	Route:      fmt.Sprintf("/%s", equates.TRANSACTIONS_TOPIC),
}

// Global DAPR client
var canxCtx context.Context
var daprclient dapr.Client
var membershipService membership.IService

func main() {
	rootCtx := context.Background()
	canxCtx, _ = signal.NotifyContext(rootCtx, os.Interrupt)

	// Create a DAPR service using a hard-coded port (must match start-membership.sh)
	s := daprd.NewService(":8085")
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

	membershipService = membership.NewService(canxCtx)
	defer membershipService.Finalize()

	// Inject dapr client and other services in actor packages
	membershipactor.Daprclient = daprclient
	membershipactor.MembershipService = membershipService

	// Register actors
	s.RegisterActorImplFactoryContext(membershipactor.MembershipActorFactory)
	fmt.Println("Membership actor registered!")

	// Register pub/sub pledge handlers
	if err := s.AddTopicEventHandler(transactionsTopicSubscription, transactionMembershipHandler); err != nil {
		panic(err)
	}
	fmt.Println("Transactions topic handler registered!")

	// Start DAPR service
	// TODO: Provide cancellation context
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func transactionMembershipHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	fmt.Println("transactionMembershipHandler....")

	go func() {
		// Decode transaction
		evt := member.MemberTransaction{}
		mapstructure.Decode(e.Data, &evt)

		fmt.Printf("Received a member transaction\n")

		// Resolve actor by membership ID
		memActorProxy := membershipactor.NewMembershipActor(evt.MembershipID)
		daprclient.ImplActorClientStub(memActorProxy)

		// Call actor methods using the main canxCtx
		err = memActorProxy.Transact(canxCtx, evt)
		if err != nil {
			fmt.Println("transactionMembershipHandler calling actor error", err)
			return
		}
	}()

	return false, nil
}

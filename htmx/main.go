package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	dapr "github.com/dapr/go-sdk/client"

	"github.com/khaledhikmat/institution-manager/htmx/flow"
	"github.com/khaledhikmat/institution-manager/htmx/server"
	"github.com/khaledhikmat/institution-manager/shared/service/campaign"
	"github.com/khaledhikmat/institution-manager/shared/service/member"
)

var stage = 0

func main() {
	rootCanx := context.Background()
	canxCtx, cancel := signal.NotifyContext(rootCanx, os.Interrupt)

	// Create a DAPR client
	// Must be a global client since it is singleton
	// Hence it would be injected in actor packages as needed
	c, err := dapr.NewClient()
	if err != nil {
		fmt.Println("Failed to start dapr client", err)
		return
	}
	defer c.Close()

	// Create service layer
	campaignService := campaign.NewService(canxCtx)
	defer campaignService.Finalize()
	memberService := member.NewService(canxCtx)
	defer memberService.Finalize()

	// Inject into workflow
	flow.DaprClient = c
	flow.CampaignService = campaignService
	flow.MemberService = memberService

	// Start workflow engine
	err = flow.Start()
	if err != nil {
		fmt.Println("Failed to start dapr workflow engine", err)
		return
	}
	defer flow.Shutdown()

	// Inject into server
	server.DaprClient = c
	server.CampaignService = campaignService
	server.MemberService = memberService

	port := "3000"
	args := os.Args[1:]
	if len(args) > 0 {
		port = args[0]
	}

	defer func() {
		cancel()
	}()

	// Launch the http server
	httpServerErr := make(chan error, 1)
	go func() {
		httpServerErr <- server.Run(canxCtx, port)
	}()

	// Wait and pump
	for {
		select {
		case err := <-httpServerErr:
			fmt.Println("error", err)
			return
		case <-canxCtx.Done():
			cancel()
			return
			// case <-time.After(20 * time.Second):
			// 	// Do something every 10 seconds
			// 	fmt.Println("Timeout")

			// 	fmt.Println("Creating a ficticious pledge....")
			// 	p := member.NewMemberPledge()
			// 	p.CampaignID = "100"
			// 	p.MemberID = "100"
			// 	p.Time = time.Now()
			// 	p.Amount = 100

			// 	// Resolve actor by campaign id
			// 	campaignActorProxy := campaignactor.NewCampaignActor(p.CampaignID)
			// 	c.ImplActorClientStub(campaignActorProxy)

			// 	// Call actor methods
			// 	err = campaignActorProxy.Pledge(canxCtx, p)
			// 	if err != nil {
			// 		fmt.Printf("Error calling the campaign actor: %v", err)
			// 	}

			// 	fmt.Println("Creating a ficticious member transaction....")
			// 	m := member.MemberTransaction{}
			// 	m.MemberID = "100"
			// 	m.MembershipID = "100"
			// 	m.Tier = "Black"
			// 	m.Type = "Purchase"
			// 	m.StartDate = time.Now()
			// 	m.ExpDate = time.Now()
			// 	m.Amount = 250
			// 	m.ExgRate = 1
			// 	m.PaymentMethod = "Visa"
			// 	m.PaymentConfirmation = "849384309483"

			// 	// Resolve actor by campaign id
			// 	memberActorProxy := memberactor.NewMemberActor(m.MemberID)
			// 	c.ImplActorClientStub(memberActorProxy)

			// 	// Call actor methods
			// 	err = memberActorProxy.Transact(canxCtx, m)
			// 	if err != nil {
			// 		fmt.Printf("Error calling the member actor: %v", err)
			// 	}
		}
	}
}

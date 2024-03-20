package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	dapr "github.com/dapr/go-sdk/client"

	"github.com/khaledhikmat/campaign-manager/htmx/server"
	campaignactor "github.com/khaledhikmat/campaign-manager/shared/actors/campaign"
	memberactor "github.com/khaledhikmat/campaign-manager/shared/actors/member"
	"github.com/khaledhikmat/campaign-manager/shared/service/campaign"
	"github.com/khaledhikmat/campaign-manager/shared/service/member"
)

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

	// Create service layer and inject into server
	server.DaprClient = c
	server.CampaignService = campaign.NewService(canxCtx)
	defer server.CampaignService.Finalize()
	server.MemberService = member.NewService(canxCtx)
	defer server.MemberService.Finalize()

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
		case <-time.After(20 * time.Second):
			// Do something every 10 seconds
			fmt.Println("Timeout")

			fmt.Println("Creating a ficticious pledge....")
			p := member.MemberPledge{}
			p.CampaignID = "100"
			p.MemberID = "100"
			p.Time = time.Now()
			p.Amount = 100

			// Resolve actor by campaign id
			campaignActorProxy := campaignactor.NewCampaignActor(p.CampaignID)
			c.ImplActorClientStub(campaignActorProxy)

			// Call actor methods
			err = campaignActorProxy.Pledge(canxCtx, p)
			if err != nil {
				fmt.Printf("Error calling the campaign actor: %v", err)
			}

			fmt.Println("Creating a ficticious member transaction....")
			m := member.MemberTransaction{}
			m.MemberID = "100"
			m.MembershipID = "100"
			m.Tier = "Black"
			m.Type = "Purchase"
			m.StartDate = time.Now()
			m.ExpDate = time.Now()
			m.Amount = 250
			m.ExgRate = 1
			m.PaymentMethod = "Visa"
			m.PaymentConfirmation = "849384309483"

			// Resolve actor by campaign id
			memberActorProxy := memberactor.NewMemberActor(m.MemberID)
			c.ImplActorClientStub(memberActorProxy)

			// Call actor methods
			err = memberActorProxy.Transact(canxCtx, m)
			if err != nil {
				fmt.Printf("Error calling the member actor: %v", err)
			}

		}
	}
}

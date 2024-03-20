package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	campaignactor "github.com/khaledhikmat/campaign-manager/shared/actors/campaign"

	"github.com/khaledhikmat/campaign-manager/shared/service/campaign"
	"github.com/khaledhikmat/campaign-manager/shared/service/member"

	dapr "github.com/dapr/go-sdk/client"
	daprd "github.com/dapr/go-sdk/service/http"
)

// Global DAPR client
var canxCtx context.Context
var daprclient dapr.Client
var cmpService campaign.IService
var memberService member.IService

func main() {
	rootCtx := context.Background()
	canxCtx, _ = signal.NotifyContext(rootCtx, os.Interrupt)

	// Create a DAPR service using a hard-coded port (must match start-campaign.sh)
	s := daprd.NewService(":8080")
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
	memberService = member.NewService(canxCtx)
	defer memberService.Finalize()

	// Inject dapr client and other services in actor packages
	campaignactor.Daprclient = daprclient
	campaignactor.CampaignService = cmpService
	campaignactor.MemberService = memberService

	// Register actors
	s.RegisterActorImplFactoryContext(campaignactor.CampaignActorFactory)
	fmt.Println("Campaign actor registered!")

	// Start DAPR service
	// TODO: Provide cancellation context
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

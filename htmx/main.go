package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/joho/godotenv"

	"github.com/khaledhikmat/institution-manager/htmx/flow"
	"github.com/khaledhikmat/institution-manager/htmx/server"
	"github.com/khaledhikmat/institution-manager/shared/service/campaign"
	"github.com/khaledhikmat/institution-manager/shared/service/member"
	"github.com/khaledhikmat/institution-manager/shared/service/realtimer"
	"github.com/khaledhikmat/institution-manager/shared/service/secret"
)

func main() {
	rootCanx := context.Background()
	canxCtx, cancel := signal.NotifyContext(rootCanx, os.Interrupt)

	// Load env vars
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Create a DAPR client
	// Must be a global client since it is singleton
	// Hence it would be injected in actor packages as needed
	c, err := dapr.NewClient()
	if err != nil {
		fmt.Println("Failed to start dapr client", err)
		return
	}
	defer c.Close()

	secretService, err := secret.NewSecretService(canxCtx, c)
	if err != nil {
		fmt.Println("Failed to create a secret service", err)
		return
	}
	defer secretService.Finalize()

	// Get Ably API Key
	secret, err := secretService.Creds("ably.apiKey")
	if err != nil {
		fmt.Println("Failed to create a secret service", err)
		//return // Do not abort for now
	}

	fmt.Println("*****Secret", secret)

	// Create service layer
	campaignService := campaign.NewService(canxCtx)
	defer campaignService.Finalize()
	memberService := member.NewService(canxCtx)
	defer memberService.Finalize()
	realtimeService, err := realtimer.NewAblyService(canxCtx, os.Getenv("ABLY_API_KEY"), "institution-manager.htmx")
	if err != nil {
		fmt.Println("Failed to start ably client", err)
		return
	}
	defer realtimeService.Finalize()

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
	server.RealtimeService = realtimeService

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

	// Wait until server exits or context is cancelled
	for {
		select {
		case err := <-httpServerErr:
			fmt.Println("http server error", err)
			return
		case <-canxCtx.Done():
			fmt.Println("application cancelled...")
			cancel()
			// Wait until downstream processors are done
			fmt.Println("wait for 2 seconds until downstream processors are cancelled...")
			time.Sleep(2 * time.Second)
			return
		}
	}
}

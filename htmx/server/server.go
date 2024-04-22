package server

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/stripe/stripe-go/v72"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/khaledhikmat/institution-manager/shared/service/campaign"
	"github.com/khaledhikmat/institution-manager/shared/service/member"
	"github.com/khaledhikmat/institution-manager/shared/service/realtimer"
)

const (
	workflowComponent = "dapr"
)

// Injected DAPR client and other services
var DaprClient dapr.Client
var CampaignService campaign.IService
var MemberService member.IService
var RealtimeService realtimer.IService

type ginWithContext func(ctx context.Context) error

func Run(canxCtx context.Context, port string) error {
	//=========================
	// ROUTER
	//=========================
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://example.com"} //TODO: Update
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "OPTIONS", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Set function map if any...

	// Link up templates and static files
	r.LoadHTMLGlob("./templates/**/*")
	r.Static("/static", "./static")

	// Load env vars
	err := godotenv.Load()
	if err != nil {
		return err
	}

	// Setup Stripe payment gateway
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Provide app info:
	stripe.SetAppInfo(&stripe.AppInfo{
		Name:    "institution-manager/payment-gateway",
		Version: "0.0.1",
		URL:     "https://github.com/khaledhikmat/institution-manager",
	})

	//=========================
	// Setup Payment ROUTES
	//=========================
	paymentRoutes(canxCtx, r)

	//=========================
	// Setup Realtime ROUTES
	//=========================
	realtimeRoutes(canxCtx, r)

	//=========================
	// Setup Home ROUTES
	//=========================
	homeRoutes(canxCtx, r)

	//=========================
	// Setup Campaign ROUTES
	//=========================
	campaignRoutes(canxCtx, r)

	//=========================
	// Setup Pledge ROUTES
	//=========================
	pledgeRoutes(canxCtx, r)

	f := cancellableGin(canxCtx, r, port)
	return f(canxCtx)
}

func getRole(_ *gin.Context) string {
	// TODO: Determine Role: Admin, Campaign Manager and Donor
	// TODO: Read from an environment variable
	return "donor"
}

func getMemberID(_ *gin.Context) string {
	// TODO: Read from an environment variable
	// Once we get the identity provider ID, we fetch from the database
	return "100"
}

func cancellableGin(_ context.Context, r *gin.Engine, port string) ginWithContext {
	return func(ctx context.Context) error {
		go func() {
			err := r.Run(":" + port)
			if err != nil {
				fmt.Println("Server start error...exiting", err)
				return
			}
		}()

		// Wait
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Server context cancelled...existing!!!")
				return ctx.Err()
			case <-time.After(time.Duration(100 * time.Second)):
				fmt.Println("Timeout....do something periodic here!!!")
			}
		}
	}
}

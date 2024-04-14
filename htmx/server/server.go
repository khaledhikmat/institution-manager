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
)

const (
	workflowComponent = "dapr"
)

// Injected DAPR client and other services
var DaprClient dapr.Client
var CampaignService campaign.IService
var MemberService member.IService

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

	// Set function map if any

	// Link up templates and static files
	r.LoadHTMLGlob("./templates/**/*")
	r.Static("/static", "./static")

	// Setup Stripe payment gateway
	err := godotenv.Load()
	if err != nil {
		return err
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Provide app info:
	stripe.SetAppInfo(&stripe.AppInfo{
		Name:    "institution-manager/payment-gateway",
		Version: "0.0.1",
		URL:     "https://github.com/khaledhikmat/institution-manager",
	})

	//=========================
	// Setup Stripe ROUTES
	//=========================
	stripeRoutes(canxCtx, r)

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

	f := cancellableGin(r, port)
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

func cancellableGin(r *gin.Engine, port string) ginWithContext {
	return func(ctx context.Context) error {
		go func() {
			r.Run(":" + port)
		}()

		// Wait
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(100 * time.Second)):
				fmt.Println("Timeout....if there is something to do!!!")
			}
		}
	}
}

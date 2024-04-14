package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/webhook"

	"github.com/gin-gonic/gin"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/khaledhikmat/institution-manager/htmx/flow"
)

const (
	workflowRefId = "workflow_ref_id"
)

// This payment payload must come from JavaScript and comply with this structure...
type payment struct {
	RefID    string `json:"refId"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

func stripeRoutes(canxCtx context.Context, r *gin.Engine) {
	//=========================
	// API
	//=========================
	r.GET("/stripe-config", func(c *gin.Context) {
		c.JSON(200, struct {
			PublishableKey string `json:"publishableKey"`
		}{
			PublishableKey: os.Getenv("STRIPE_PUBLISHABLE_KEY"),
		})
	})

	r.POST("/stripe-webhook", func(c *gin.Context) {
		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, struct {
				Error string `json:"error"`
			}{
				Error: err.Error(),
			})
			return
		}

		event, err := webhook.ConstructEvent(b, c.GetHeader("Stripe-Signature"), os.Getenv("STRIPE_WEBHOOK_SECRET"))
		if err != nil {
			c.JSON(500, struct {
				Error string `json:"error"`
			}{
				Error: err.Error(),
			})
			return
		}

		if event.Type == "checkout.session.completed" {
			fmt.Println("Checkout Session completed!")
		}

		if event.Type == /*stripe.EventTypePaymentIntentSucceeded*/ "payment_intent.succeeded" {
			fmt.Println("💰 Payment succeeded!")
			// prettyJSON, _ := json.MarshalIndent(event, "", "    ")
			// fmt.Printf("EVENT: %s\n", prettyJSON)

			// Read metdata to determine the workflow ID
			workflowID := ""
			metadata := event.Data.Object["metadata"]
			if metadata != nil {
				v, ok := metadata.(map[string]interface{})
				if ok {
					workflowID = fmt.Sprintf("%v", v[workflowRefId])
				}
			}

			if workflowID == "" {
				// TODO: Report an error
				fmt.Println("😢 Received a webhook without a valid workflow")
			}

			// Using the reference id, trigger a pay event
			err := DaprClient.RaiseEventWorkflowBeta1(canxCtx, &dapr.RaiseEventWorkflowRequest{
				InstanceID:        workflowID,
				WorkflowComponent: workflowComponent,
				EventName:         flow.PledgePaidEvent,
				EventData:         nil,
				SendRawData:       false,
			})

			if err != nil {
				// TODO: Report an error
				fmt.Printf("😢 Received an error trying to raise an event on workflow %s: %v\n", workflowID, err)
			}
		}

		c.JSON(200, struct{}{})
	})

	r.POST("/stripe-payment-intent", func(c *gin.Context) {
		var p payment
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(500, struct {
				Error string `json:"error"`
			}{
				Error: err.Error(),
			})
			return
		}

		amt, err := strconv.ParseInt(p.Amount, 10, 64)
		if err != nil {
			c.JSON(500, struct {
				Error string `json:"error"`
			}{
				Error: err.Error(),
			})
			return
		}

		params := &stripe.PaymentIntentParams{
			Amount: stripe.Int64(amt),
			// Currency: stripe.String(string(stripe.CurrencyUSD)),
			Currency: stripe.String(p.Currency),
			AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
				Enabled: stripe.Bool(true),
			},
		}
		params.AddMetadata(workflowRefId, p.RefID)

		pi, err := paymentintent.New(params)
		if err != nil {
			if stripeErr, ok := err.(*stripe.Error); ok {
				c.JSON(500, struct {
					Error string `json:"error"`
				}{
					Error: "stripe error " + stripeErr.Error(),
				})
				return
			}

			c.JSON(500, struct {
				Error string `json:"error"`
			}{
				Error: err.Error(),
			})
			return
		}

		c.JSON(200, struct {
			ClientSecret string `json:"clientSecret"`
		}{
			ClientSecret: pi.ClientSecret,
		})
	})
}

package server

import (
	"context"
	"fmt"
	"strconv"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/gin-gonic/gin"
	"github.com/khaledhikmat/institution-manager/htmx/flow"
	campaignactor "github.com/khaledhikmat/institution-manager/shared/actors/campaign"
	"github.com/khaledhikmat/institution-manager/shared/service/campaign"
	"github.com/khaledhikmat/institution-manager/shared/service/member"
)

func pledgeRoutes(canxCtx context.Context, r *gin.Engine) {
	//=========================
	// PAGES
	//=========================
	r.GET("/pledges", func(c *gin.Context) {
		role := getRole(c)
		target := "pledges.html"
		if c.Query("id") == "" {
			c.HTML(200, target, gin.H{
				"Tab":   "Pledges",
				"Role":  role,
				"Error": "Campaign id is missing!",
			})
			return
		}

		if c.GetHeader("HX-Request") == "true" {
			target = "pledges-list.html"
		}

		pledges, err := MemberService.GetPledgesByCampaign(c.Query("id"))
		if err != nil {
			c.HTML(200, target, gin.H{
				"Tab":     "Pledges",
				"Role":    role,
				"Error":   err.Error(),
				"Pledges": pledges,
			})
			return
		}

		c.HTML(200, target, gin.H{
			"Tab":     "Pledges",
			"Role":    role,
			"Error":   "",
			"Pledges": pledges,
		})
	})

	r.GET("/pledge", func(c *gin.Context) {
		role := getRole(c)
		target := "pledge.html"
		if c.Query("id") == "" {
			c.HTML(200, target, gin.H{
				"Tab":   "Home",
				"Role":  role,
				"Error": "Campaign id is missing!",
			})
			return
		}

		campaign, err := CampaignService.GetCampaign(c.Query("id"))
		if err != nil {
			c.HTML(200, target, gin.H{
				"Tab":      "Home",
				"Role":     role,
				"Error":    err.Error(),
				"Campaign": campaign,
			})
			return
		}

		c.HTML(200, target, gin.H{
			"Tab":      "Home",
			"Role":     role,
			"Error":    "",
			"Campaign": campaign,
		})
	})

	r.GET("/pledgepayment", func(c *gin.Context) {
		target := "pledge-payment.html"
		if c.Query("i") == "" {
			target = "pledge-response.html"
			c.HTML(200, target, gin.H{
				"Mesage": "reference id is missing!",
			})
			return
		}

		if c.Query("a") == "" {
			target = "pledge-response.html"
			c.HTML(200, target, gin.H{
				"Message": "amount is missing!",
			})
			return
		}

		amt, err := strconv.ParseInt(c.Query("a"), 10, 64)
		if err != nil {
			target = "pledge-response.html"
			c.HTML(200, target, gin.H{
				"Message": fmt.Sprintf("amount is not parseable - %v", err),
			})
			return
		}

		amt = amt * 100 // To satisfy Stripe

		if c.Query("c") == "" {
			target = "pledge-response.html"
			c.HTML(200, target, gin.H{
				"Message": "currency is missing!",
			})
			return
		}

		c.HTML(200, target, gin.H{
			"Message":   "",
			"Reference": c.Query("i"),
			"Amount":    amt,
			"Currency":  c.Query("c"),
		})
	})

	//=========================
	// ACTIONS
	//=========================
	// Post a pledge
	r.POST("/actions/pledges", func(c *gin.Context) {
		memberID := getMemberID(c)
		target := "pledge-response.html"

		if c.PostForm("id") == "" {
			c.HTML(200, target, gin.H{
				"Message": "Campaign id is missing!",
			})
			return
		}

		// Validate campaign
		camp, err := CampaignService.GetCampaign(c.PostForm("id"))
		if err != nil {
			c.HTML(200, target, gin.H{
				"Message": err.Error(),
			})
			return
		}

		// Apply form data
		p := member.NewMemberPledge()
		p.MemberID = memberID
		p.CampaignID = camp.ID
		p.Time = time.Now()
		amt, err := strconv.ParseInt(c.PostForm("amount"), 10, 64)
		if err != nil {
			c.HTML(200, target, gin.H{
				"Message": fmt.Sprintf("Unable to pledge %v", err),
			})
			return
		}
		p.Amount = amt

		// If campaign type is confirmed and payment type is pre-pledge:
		// 1. Launch a pre-payment pledge workflow
		// 2. Re-direct user to payment page
		// 3. Once the payment is completed (via webhook), the workflow 'paid' event will be triggered to submit the pledge.
		if camp.Type == campaign.ConfirmedCampaign && camp.PaymentType == campaign.PrePledgePaymentCampaign {
			respPrePay, err := DaprClient.StartWorkflowBeta1(canxCtx, &dapr.StartWorkflowRequest{
				InstanceID:        p.ID,
				WorkflowComponent: workflowComponent,
				WorkflowName:      flow.PrePledgePaymentWorkflowName,
				Options:           nil,
				Input: flow.CampaignPledge{
					Pledge:   p,
					Campaign: camp,
				},
				SendRawInput: false,
			})
			if err != nil {
				c.HTML(200, target, gin.H{
					"Message": fmt.Sprintf("Unable to start pre-payment pledge workflow %v", err),
				})
				return
			}

			fmt.Printf("pre-payment pledge workflow started with id: %v\n", respPrePay.InstanceID)

			// Redirect user to pay
			c.Redirect(303, fmt.Sprintf("/pledgepayment?i=%s&a=%d&c=%s", p.ID, p.Amount, camp.Currency))
			return
		}

		// Resolve actor by campaign id
		campaignActorProxy := campaignactor.NewCampaignActor(p.CampaignID)
		DaprClient.ImplActorClientStub(campaignActorProxy)

		// Call actor methods
		err = campaignActorProxy.Pledge(canxCtx, p)
		if err != nil {
			c.HTML(200, target, gin.H{
				"Message": fmt.Sprintf("Unable to pledge %v", err),
			})
			return
		}

		// If campaign type is confirmed and paymet type is post-pledge:
		// 1. Launch a pre-payment pledge workflow
		// 2. Re-direct user to payment page
		// 3. Once the payment is completed (via webhook), the workflow 'paid' event will be triggered to confirm the pledge.
		if camp.Type == campaign.ConfirmedCampaign && camp.PaymentType == campaign.PostPledgePaymentCampaign {
			respPostPay, err := DaprClient.StartWorkflowBeta1(canxCtx, &dapr.StartWorkflowRequest{
				InstanceID:        p.ID,
				WorkflowComponent: workflowComponent,
				WorkflowName:      flow.PostPledgePaymentWorkflowName,
				Options:           nil,
				Input: flow.CampaignPledge{
					Pledge:   p,
					Campaign: camp,
				},
				SendRawInput: false,
			})
			if err != nil {
				c.HTML(200, target, gin.H{
					"Message": fmt.Sprintf("Unable to start post-payment pledge workflow %v", err),
				})
				return
			}

			fmt.Printf("post-payment pledge workflow started with id: %v\n", respPostPay.InstanceID)

			// Redirect user to pay
			c.Redirect(303, fmt.Sprintf("/pledgepayment?i=%s&a=%d&c=%s", p.ID, p.Amount, camp.Currency))
			return
		}

		c.HTML(200, target, gin.H{
			"Message": "Thank you!",
		})
	})

	// Post a payment - the callback from Stripe is what we need to do to complete the payment
	r.POST("/actions/pledgepayments", func(c *gin.Context) {
		target := "payment-response.html"

		// TODO: To allow the form to submit!!!
		time.Sleep(3 * time.Second)

		c.HTML(200, target, gin.H{
			"Message": "Thank you for your payment! We will confirm via email.",
		})
	})
}

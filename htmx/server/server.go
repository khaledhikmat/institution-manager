package server

import (
	"context"
	"fmt"
	"strconv"
	"time"

	dapr "github.com/dapr/go-sdk/client"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	campaignactor "github.com/khaledhikmat/campaign-manager/shared/actors/campaign"
	"github.com/khaledhikmat/campaign-manager/shared/service/campaign"
	"github.com/khaledhikmat/campaign-manager/shared/service/member"
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

	r.LoadHTMLGlob("./htmx/templates/**/*")
	r.Static("/static", "./static")

	//=========================
	// PAGES
	//=========================
	r.GET("/", func(c *gin.Context) {
		// TODO: Determine Role: Admin, Campaign Manager and Donor
		role := "donor"
		target := "index.html"
		if c.GetHeader("HX-Request") == "true" {
			target = "campaigns-list.html"
		}

		campaigns, err := CampaignService.GetCampaigns("")
		if err != nil {
			c.HTML(200, target, gin.H{
				"Tab":       "Home",
				"Role":      role,
				"Error":     err.Error(),
				"Campaigns": campaigns,
			})
			return
		}

		c.HTML(200, target, gin.H{
			"Tab":       "Home",
			"Role":      role,
			"Error":     "",
			"Campaigns": campaigns,
		})
	})

	r.GET("/campaign", func(c *gin.Context) {
		// TODO: Determine Role: Admin, Campaign Manager and Donor
		role := "donor"
		target := "campaign.html"
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

	r.GET("/pledges", func(c *gin.Context) {
		// TODO: Determine Role: Admin, Campaign Manager and Donor
		role := "donor"
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
		// TODO: Determine Role: Admin, Campaign Manager and Donor
		role := "donor"
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

	//=========================
	// ACTIONS
	//=========================
	r.GET("/actions/load-more-campaigns", func(c *gin.Context) {
		c.Redirect(303, fmt.Sprintf("/?t=%s", c.Query("t")))
	})

	// New Campaign - add a new item in the database, create and initialize an actor
	r.POST("/actions/campaigns", func(c *gin.Context) {
		// TODO: Determine Role: Admin, Campaign Manager and Donor
		role := "donor"
		target := "new-campaign.html"

		if role != "admin" {
			c.HTML(200, target, gin.H{
				"Tab":   "Home",
				"Role":  role,
				"Error": "restricted operation",
			})
			return
		}

		cmp := campaign.Campaign{}
		// Apply form data
		campaignFormData(c, &cmp)

		// Persist to the database
		_, err := CampaignService.NewCampaign(cmp)
		if err != nil {
			c.HTML(200, target, gin.H{
				"Tab":   "Home",
				"Role":  role,
				"Error": fmt.Sprintf("Unable to create campaign %v", err),
			})
			return
		}

		c.HTML(200, target, gin.H{
			"Tab":   "Home",
			"Role":  role,
			"Error": "",
		})
	})

	// Update Campaign - refresh actor and let the actor update the database
	r.PUT("/actions/campaigns", func(c *gin.Context) {
		// TODO: Determine Role: Admin, Campaign Manager and Donor
		role := "admin"
		target := "edit-campaign.html"

		if role != "admin" {
			c.HTML(200, target, gin.H{
				"Tab":   "Home",
				"Role":  role,
				"Error": "restricted operation",
			})
			return
		}

		id := c.PostForm("id")
		cmp, err := CampaignService.GetCampaign(id)
		if err != nil {
			c.HTML(200, target, gin.H{
				"Tab":   "Home",
				"Role":  role,
				"Error": fmt.Sprintf("Unable to retrieve campaign %s", id),
			})
			return
		}

		// Apply form data
		campaignFormData(c, &cmp)

		// Resolve actor by campaign id
		campaignActorProxy := campaignactor.NewCampaignActor(cmp.Id)
		DaprClient.ImplActorClientStub(campaignActorProxy)

		// Call actor methods
		err = campaignActorProxy.Update(canxCtx, cmp)
		if err != nil {
			c.HTML(200, target, gin.H{
				"Tab":   "Home",
				"Role":  role,
				"Error": fmt.Sprintf("Unable to update campaign %v", err),
			})
			return
		}

		c.HTML(200, target, gin.H{
			"Tab":   "Home",
			"Role":  role,
			"Error": "",
		})
	})

	// Post a pledge
	r.POST("/actions/pledges", func(c *gin.Context) {
		target := "pledge-response.html"

		if c.PostForm("id") == "" {
			c.HTML(200, target, gin.H{
				"Message": "Campaign id is missing!",
			})
			return
		}

		_, err := CampaignService.GetCampaign(c.PostForm("id"))
		if err != nil {
			c.HTML(200, target, gin.H{
				"Message": err.Error(),
			})
			return
		}

		// Apply form data
		p := member.MemberPledge{}
		p.CampaignID = c.PostForm("id")
		p.Time = time.Now()
		amt, err := strconv.ParseInt(c.PostForm("amount"), 10, 64)
		if err != nil {
			c.HTML(200, target, gin.H{
				"Message": fmt.Sprintf("Unable to pledge %v", err),
			})
			return
		}
		p.Amount = amt

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

		c.HTML(200, target, gin.H{
			"Message": "Thank you!",
		})
	})

	f := cancellableGin(r, port)
	return f(canxCtx)
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

func campaignFormData(c *gin.Context, cmp *campaign.Campaign) error {
	// Apply form data
	cmp.Name = c.PostForm("name")
	cmp.Description = c.PostForm("desc")
	v, err := strconv.ParseBool(c.PostForm("virtual"))
	if err != nil {
		return err
	}
	cmp.Virtual = v
	cmp.ImageUrl = c.PostForm("imageurl")
	g, err := strconv.ParseInt(c.PostForm("goal"), 10, 64)
	if err != nil {
		return err
	}
	cmp.Goal = g
	cmp.Currency = c.PostForm("currency")

	d, err := strconv.ParseInt(c.PostForm("duration"), 10, 64)
	if err != nil {
		return err
	}
	cmp.Duration = d
	return nil
}

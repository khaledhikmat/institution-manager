package server

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	campaignactor "github.com/khaledhikmat/institution-manager/shared/actors/campaign"
	"github.com/khaledhikmat/institution-manager/shared/service/campaign"
)

func campaignRoutes(canxCtx context.Context, r *gin.Engine) {
	//=========================
	// PAGES
	//=========================
	r.GET("/campaigns", func(c *gin.Context) {
		role := getRole(c)
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
		role := getRole(c)
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

	//=========================
	// ACTIONS
	//=========================
	r.GET("/actions/load-more-campaigns", func(c *gin.Context) {
		c.Redirect(303, fmt.Sprintf("/?t=%s", c.Query("t")))
	})

	// New Campaign - add a new item in the database, create and initialize an actor
	r.POST("/actions/campaigns", func(c *gin.Context) {
		role := getRole(c)
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
		role := getRole(c)
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
		campaignActorProxy := campaignactor.NewCampaignActor(cmp.ID)
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
	cmp.ImageURL = c.PostForm("imageurl")
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

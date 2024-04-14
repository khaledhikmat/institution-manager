package server

import (
	"context"

	"github.com/gin-gonic/gin"
)

func homeRoutes(_ context.Context, r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
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
}

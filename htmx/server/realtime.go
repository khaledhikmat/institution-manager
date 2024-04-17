package server

import (
	"context"

	"github.com/gin-gonic/gin"
)

func realtimeRoutes(_ context.Context, r *gin.Engine) {
	//=========================
	// API
	//=========================
	r.POST("/realtime-auth", func(c *gin.Context) {
		// WARNING: We are assuming that the user is authenticated because they are accessing HTMX
		token, err := RealtimeService.Token()
		if err != nil {
			c.JSON(500, "")
			return
		}

		c.JSON(200, token)
	})
}

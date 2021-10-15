package server

import (
	"net/http"

	"tldw-server/internal/api"

	"github.com/gin-gonic/gin"
)

func registerRoutes(router *gin.Engine) {
	router.RedirectTrailingSlash = true

	v1 := router.Group("/api/v1")
	{
		// Health
		api.GetStatus(v1)

		// Main Logic
		api.GetSummary(v1)
		api.GetSummaryByID(v1)
		api.SubmitSummary(v1)
		api.SummaryVote(v1)

		// User
		api.SignupUser(v1)
	}

	// Default HTML page
	router.NoRoute(
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		},
	)
}
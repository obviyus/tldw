package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"tldw-server/internal/models"
)

type CreateUser struct {
	Username string `json:"username"`
	Platform string `json:"platform"`
}

// SignupUser POST /api/v1/signup
// Create a new user given a username and return an API key
//
// Parameters:
// - username: unique username for user
// - platform: which platform user is on
func SignupUser(router *gin.RouterGroup) {
	router.POST(
		"/signup", func(c *gin.Context) {
			var params CreateUser
			if err := c.BindJSON(&params); err == nil {
				// Check if username exists
				if models.FindUserByName(params.Username) != nil {
					AbortAlreadyExists(c, fmt.Sprintf("Username %s", params.Username))
					return
				}

				if user := models.NewUser(params.Username, params.Platform); user != nil {
					if err := user.Create(); err != nil {
						AbortUnexpected(c)
						return
					}

					apiKey := models.NewApiKey(user)
					if apiKey == "" {
						AbortUnexpected(c)
						return
					}

					c.JSON(
						http.StatusOK, gin.H{
							"user": user,
							"key":  apiKey,
						},
					)
				} else {
					AbortUnexpected(c)
				}
			} else {
				AbortBadRequest(c)
			}
		},
	)
}

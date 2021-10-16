package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"tldw-server/internal/models"
	"tldw-server/internal/tldw"
)

// GetSummary GET /api/v1/summary/:videoId
//
// Parameters:
// - videoId: ID of the target video as a path parameter
func GetSummary(router *gin.RouterGroup) {
	router.GET(
		"/summary/:videoId", func(c *gin.Context) {
			videoId := c.Param("videoId")
			if videoId == "" {
				AbortMissingParameter(c, "missing parameter: videoId")
				return
			} else if !tldw.IsValidSummary(videoId) {
				AbortBadRequest(c)
				return
			}

			// Check if videoID exists in database
			if summaries := models.FindSummariesForVideo(videoId); len(summaries) != 0 {
				c.JSON(
					http.StatusOK, gin.H{
						"summaries": summaries,
					},
				)
			} else {
				AbortEntityNotFound(c)
			}
		},
	)
}

// GetSummaryByID GET /api/v1/summary/:videoId/:id
//
// Parameters:
// - videoId: ID of the target video as a path parameter
// - Id: ID of the summary entity
func GetSummaryByID(router *gin.RouterGroup) {
	router.GET(
		"/summary/:videoId/:Id", func(c *gin.Context) {
			videoId, Id := c.Param("videoId"), c.Param("Id")
			if videoId == "" || Id == "" {
				AbortMissingParameter(c, "missing parameter in path")
				return
			} else if !tldw.IsValidVideoID(videoId) {
				AbortBadRequest(c)
				return
			}

			// Check if Id exists in database
			if summary := models.FindSummaryByID(Id); summary != nil {
				c.JSON(
					http.StatusOK, gin.H{
						"summary": summary,
					},
				)
			} else {
				AbortEntityNotFound(c)
			}
		},
	)
}

type summaryParams struct {
	Summary  string `json:"summary" binding:"required"`
	Language string `json:"language"`
	ApiKey   string `json:"key" binding:"required"`
}

// SubmitSummary POST /api/v1/summary/:videoId
//
// Parameters
// - summary: summary of video
// - language: language of summary
// - key: unique ApiKey
func SubmitSummary(router *gin.RouterGroup) {
	router.POST(
		"/summary/:videoId", func(c *gin.Context) {
			videoId := c.Param("videoId")
			if videoId == "" {
				AbortMissingParameter(c, "missing parameter: videoId")
				return
			} else if !tldw.IsValidVideoID(videoId) {
				log.Error("invalid videoId: ", videoId)
				AbortBadRequest(c)
				return
			}

			var params summaryParams
			if err := c.BindJSON(&params); err == nil {
				// Validate API key
				key := models.FindKey(params.ApiKey)
				if key == nil {
					AbortUnauthorized(c)
					return
				}

				// Check if user exists
				if user := models.FindUserByUserID(key.ID); user == nil {
					AbortUnauthorized(c)
					return
				} else {
					// Validate User Input
					if !tldw.IsValidSummary(params.Summary) {
						AbortBadRequest(c)
						return
					} else if tldw.CheckProfanity(params.Summary) {
						AbortProfanity(c)
						return
					}

					summary, err := user.AddSummary(params.Summary, videoId, params.Language)
					if err == nil {
						c.JSON(
							http.StatusOK, gin.H{
								"summary": summary,
							},
						)
					} else {
						AbortUnexpected(c)
						log.Fatalln(err)
					}
				}
			} else {
				log.Error(err)
				AbortBadRequest(c)
			}
		},
	)
}

type VoteParams struct {
	Modifier bool   `json:"modifier"`
	ApiKey   string `json:"key"`
}

// SubmitVote POST /api/v1/summary/:videoId/:Id
//
// Parameters:
// modifier: boolean value to either increment or decrement
func SubmitVote(router *gin.RouterGroup) {
	router.POST(
		"/summary/:videoId/:Id", func(c *gin.Context) {
			var params VoteParams
			videoId, Id := c.Param("videoId"), c.Param("Id")
			if videoId == "" || Id == "" {
				AbortMissingParameter(c, "missing parameter in path")
				return
			} else if !tldw.IsValidVideoID(videoId) {
				AbortBadRequest(c)
				return
			} else if err := c.BindJSON(&params); err != nil {
				AbortBadRequest(c)
				return
			}

			// Validate API key
			if key := models.FindKey(params.ApiKey); key == nil {
				AbortUnauthorized(c)
				return
			} else {
				// Check if user exists
				if user := models.FindUserByUserID(key.ID); user == nil {
					AbortUnauthorized(c)
					return
				} else {
					// Check if SummaryID exists in database
					if summary := models.FindSummaryByID(Id); summary == nil {
						AbortEntityNotFound(c)
					} else {
						var value int
						if params.Modifier {
							value = 1
						} else {
							value = -1
						}

						// Check if user has voted for this Summary
						if vote := summary.SummaryUserVote(user); vote == nil {
							if newVote := summary.SubmitVote(user.ID, value); newVote != nil {
								c.JSON(
									http.StatusOK,
									gin.H{
										"vote": newVote,
									},
								)
							} else {
								AbortUnexpected(c)
							}
						} else {
							if vote.Value == value {
								AbortAlreadyExists(c, fmt.Sprintf("vote %s", vote.ID))
							} else {
								if newVote := vote.UpdateVote(value); newVote == nil {
									c.JSON(
										http.StatusOK,
										gin.H{
											"vote": newVote,
										},
									)
								} else {
									AbortUnexpected(c)
								}
							}
						}
					}
				}
			}
		},
	)
}

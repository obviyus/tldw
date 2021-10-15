package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tldw-server/internal/event"
	"tldw-server/internal/i18n"
)

var log = event.Log

// Abort aborts any HTTP status codes with a JSON response
func Abort(c *gin.Context, code int, id i18n.Message, params ...interface{}) {
	resp := i18n.NewResponse(code, id, params...)

	log.Debugf(
		"api: aborted %s with code %d: (%s)", c.FullPath(), code, resp.String(),
	)
	c.AbortWithStatusJSON(code, resp)
}

// AbortSaveFailed handles aborting for HTTP 500 codes
func AbortSaveFailed(c *gin.Context) {
	Abort(c, http.StatusInternalServerError, i18n.ErrSaveFailed)
}

// AbortDeleteFailed handles aborting for HTTP 500 codes
func AbortDeleteFailed(c *gin.Context) {
	Abort(c, http.StatusInternalServerError, i18n.ErrDeleteFailed)
}

// AbortUnexpected handles aborting for HTTP 500 codes
func AbortUnexpected(c *gin.Context) {
	Abort(c, http.StatusInternalServerError, i18n.ErrUnexpected)
}

// AbortUnauthorized handles aborting for HTTP 401 codes
func AbortUnauthorized(c *gin.Context) {
	Abort(c, http.StatusUnauthorized, i18n.ErrUnauthorized)
}

// AbortEntityNotFound handles aborting for HTTP 404 codes
func AbortEntityNotFound(c *gin.Context) {
	Abort(c, http.StatusNotFound, i18n.ErrEntityNotFound)
}

// AbortBadRequest handles aborting for HTTP 400 codes
func AbortBadRequest(c *gin.Context) {
	Abort(c, http.StatusBadRequest, i18n.ErrBadRequest)
}

// AbortAlreadyExists handles aborting for HTTP 409 codes
func AbortAlreadyExists(c *gin.Context, s string) {
	Abort(c, http.StatusConflict, i18n.ErrAlreadyExists, s)
}

// AbortMissingParameter handles aborting for HTTP 400 codes
func AbortMissingParameter(c *gin.Context, s string) {
	Abort(c, http.StatusBadRequest, i18n.ErrMissingParameter, s)
}

// AbortProfanity handles aborting with profanity in messages
func AbortProfanity(c *gin.Context) {
	Abort(c, http.StatusBadRequest, i18n.ErrProfanity)
}

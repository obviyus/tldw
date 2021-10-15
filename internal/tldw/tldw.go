package tldw

import (
	"log"
	"regexp"
	"unicode/utf8"

	"github.com/TwiN/go-away"
	"github.com/microcosm-cc/bluemonday"
)

func IsValidVideoID(videoID string) bool {
	// https://webapps.stackexchange.com/a/101153
	matched, err := regexp.MatchString("[0-9A-Za-z_-]{10}[048AEIMQUYcgkosw]", videoID)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return matched
}

func CheckProfanity(input string) bool {
	return goaway.IsProfane(input)
}

func IsValidSummary(summary string) bool {
	p := bluemonday.UGCPolicy()
	return utf8.RuneCountInString(summary) < 280 && p.Sanitize(summary) == summary
}

package models

import (
	"sync"

	"gorm.io/gorm"

	"tldw-server/pkg/rnd"
	"tldw-server/pkg/txt"
)

var SummaryMutex = sync.Mutex{}

type Summaries []Summary

type Summary struct {
	gorm.Model
	ID       string
	VideoID  string
	UserID   string
	Summary  string
	Score    int
	Language string
}

func (Summary) TableName() string {
	return "summaries"
}

// Create : inserts a new row to the database
func (s *Summary) Create() error {
	SummaryMutex.Lock()
	defer SummaryMutex.Unlock()

	return g.Db().Create(s).Error
}

// Delete removes the summary from the database.
func (s *Summary) Delete() error {
	SummaryMutex.Lock()
	defer SummaryMutex.Unlock()

	return g.Db().Delete(s).Error
}

// Save : saves the new row to the database
func (s *Summary) Save() error {
	SummaryMutex.Lock()
	defer SummaryMutex.Unlock()

	return g.Db().Save(s).Error
}

// FindSummariesForVideo returns an existing Summary or nil if not found.
func FindSummariesForVideo(videoID string) (result []Summary) {
	if err := g.Db().Where(
		"video_id = ? AND score > -5", videoID,
	).Order("score desc").Limit(5).Find(&result).Error; err == nil {
		for _, summary := range result {
			summary.Score = FindVotesForSummary(summary.ID)
		}

		return result
	} else {
		log.Debugf("user %s not found", txt.Quote(videoID))
		return nil
	}
}

// FindSummaryByID returns an existing Summary or nil if not found.
func FindSummaryByID(SummaryID string) (result *Summary) {
	if err := g.Db().Where(
		"id = ?", SummaryID,
	).First(&result).Error; err == nil {
		owner := FindUserByUserID(result.UserID)
		if err := g.Db().Model(owner).UpdateColumn(
			"requests_made", gorm.Expr("requests_made + ?", 1),
		).Error; err != nil {
			log.Errorf("user: %s (update number of requests)", err)
		}

		result.Score = FindVotesForSummary(result.ID)
		if err != nil {
			log.Error(err)
		}
		return result
	} else {
		log.Debugf("user %s not found", txt.Quote(SummaryID))
		return nil
	}
}

// SummaryUserVote returns a Vote if User has already voted else nil
func (s *Summary) SummaryUserVote(u *User) (result *Vote) {
	if err := g.Db().Model(&Vote{}).Where("summary_id = ? AND user_id = ?", s.ID, u.ID).First(&result).Error; err == nil {
		return result
	} else {
		return nil
	}
}

// SubmitVote adds a vote for a user to a summaryID
func (s *Summary) SubmitVote(userID string, value int) (result *Vote) {
	newVote := Vote{
		ID:        rnd.TLDWID('v'),
		SummaryID: s.ID,
		UserID:    userID,
		Value:     value,
	}

	if err := newVote.Create(); err == nil {
		return &newVote
	} else {
		log.Error(err)
		return nil
	}
}

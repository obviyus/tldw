package models

import (
	"sync"

	"gorm.io/gorm"

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

		return result
	} else {
		log.Debugf("user %s not found", txt.Quote(SummaryID))
		return nil
	}
}

// UpdateScore modifies the summary score of an entity
func (s *Summary) UpdateScore(modifier int) (result *Summary) {
	s.Score += modifier
	if err := s.Save(); err != nil {
		return nil
	}

	return s
}

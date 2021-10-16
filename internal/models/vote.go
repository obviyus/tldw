package models

import (
	"sync"

	"gorm.io/gorm"

	"tldw-server/pkg/txt"
)

var VoteMutex = sync.Mutex{}

type Votes []Vote

type Vote struct {
	gorm.Model
	ID        string
	SummaryID string
	UserID    string
	Value     int
}

func (Vote) TableName() string {
	return "votes"
}

// Create : inserts a new row to the database
func (v *Vote) Create() error {
	VoteMutex.Lock()
	defer VoteMutex.Unlock()

	return g.Db().Create(v).Error
}

// Delete removes the summary from the database.
func (v *Vote) Delete() error {
	VoteMutex.Lock()
	defer VoteMutex.Unlock()

	return g.Db().Delete(v).Error
}

// Save : saves the new row to the database
func (v *Vote) Save() error {
	VoteMutex.Lock()
	defer VoteMutex.Unlock()

	return g.Db().Save(v).Error
}

type VoteSum struct {
	total int
}

// FindVotesForSummary returns an existing Vote or nil if not found.
func FindVotesForSummary(summaryID string) (total *VoteSum) {
	if err := g.Db().Table("votes").Select("sum(value) as total").Where("summary_id = ?", summaryID).First(&total); err == nil {
		return total
	} else {
		log.Error(err)
	}

	return &VoteSum{total: 0}
}

// FindVotesForUser returns an existing Vote or nil if not found.
func FindVotesForUser(userID string) (total *VoteSum) {
	if err := g.Db().Table("votes").Select("sum(value) as total").Where("user_id = ?", userID).First(&total); err == nil {
		return total
	} else {
		log.Error(err)
	}

	return &VoteSum{total: 0}
}

// FindVoteByID returns an existing Vote or nil if not found.
func FindVoteByID(VoteID string) (result *Vote) {
	if err := g.Db().Where(
		"id = ?", VoteID,
	).First(&result).Error; err == nil {
		owner := FindUserByUserID(result.UserID)
		if err := g.Db().Model(owner).UpdateColumn(
			"requests_made", gorm.Expr("requests_made + ?", 1),
		).Error; err != nil {
			log.Errorf("user: %s (update number of requests)", err)
		}

		return result
	} else {
		log.Debugf("user %s not found", txt.Quote(VoteID))
		return nil
	}
}

// UpdateVote modifies value of a vote
func (v *Vote) UpdateVote(value int) (result *Vote) {
	v.Value = value
	if err := v.Save(); err != nil {
		return nil
	}

	return v
}

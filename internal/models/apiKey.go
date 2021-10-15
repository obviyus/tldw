package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// ApiKey is unique and tied to each user
type ApiKey struct {
	gorm.Model
	ID  string
	Key string
}

func (ApiKey) TableName() string {
	return "api_keys"
}

// NewApiKey creates a new ApiKey object and stores it in the database
func NewApiKey(user *User) (key string) {
	if user == nil || user.ID == "" {
		return ""
	}

	key = uuid.NewV4().String()

	a := ApiKey{
		ID:  user.ID,
		Key: key,
	}

	if err := a.Save(); err != nil {
		log.Fatal(err)
		return ""
	}

	return key
}

// Create inserts a new row to the database.
func (key *ApiKey) Create() error {
	return g.Db().Create(key).Error
}

// Save inserts a new row to the database or updates a row if the primary key already exists.
func (key *ApiKey) Save() error {
	return g.Db().Save(key).Error
}

// FindKey returns an ApiKey object if given key exists in database.
func FindKey(key string) *ApiKey {
	result := ApiKey{}
	if err := g.Db().Where("key = ?", key).First(&result).Error; err == nil {
		return &result
	}

	return nil
}

package models

import (
	"tldw-server/pkg/rnd"
	"tldw-server/pkg/txt"

	"gorm.io/gorm"
)

type Users []User

type User struct {
	gorm.Model
	ID           string
	UserName     string
	RequestsMade int
	Summaries    []Summary
	Votes        int
	Platform     string
}

// NewUser creates a user entity
func NewUser(username string, platform string) *User {
	if username == "" {
		return nil
	}

	return &User{
		ID:           rnd.TLDWID('u'),
		UserName:     username,
		RequestsMade: 0,
		Votes:        0,
		Platform:     platform,
	}
}

func (User) TableName() string {
	return "users"
}

// Admin : Default admin user.
var Admin = User{
	ID:       "u000000000000000",
	UserName: "admin",
}

// UnknownUser : Anonymous, public user without own account.
var UnknownUser = User{
	ID: "u000000000000001",
}

var GuestUser = User{
	ID:       "u000000000000002",
	UserName: "Guest",
}

// CreateDefaultUsers initializes the database with default user accounts.
func CreateDefaultUsers() {
	if user := FirstOrCreateUser(&Admin); user != nil {
		Admin = *user
	}

	if user := FirstOrCreateUser(&UnknownUser); user != nil {
		UnknownUser = *user
	}

	if user := FirstOrCreateUser(&GuestUser); user != nil {
		GuestUser = *user
	}
}

// Create : inserts a new row to the database
func (u *User) Create() error {
	return g.Db().Create(u).Error
}

// Save : saves the new row to the database
func (u *User) Save() error {
	return g.Db().Save(u).Error
}

// FirstOrCreateUser : returns a row if it exists, else inserts a new row and returns that
func FirstOrCreateUser(u *User) (result *User) {
	g.Db().FirstOrCreate(&result, u)

	return result
}

// FindUserByUserID returns an existing user or nil if not found.
func FindUserByUserID(userID string) (result *User) {
	if err := g.Db().Where("id = ?", userID).First(&result).Error; err == nil {
		if err := g.Db().Model(result).UpdateColumn(
			"requests_made", gorm.Expr("requests_made + ?", 1),
		).Error; err != nil {
			log.Errorf("user: %s (update login attempts)", err)
		}
		return result
	} else {
		log.Debugf("user %s not found", txt.Quote(userID))
		return nil
	}
}

// FindUserByName returns an existing user or nil if not found.
func FindUserByName(userName string) *User {
	if userName == "" {
		return nil
	}

	result := User{}
	if err := g.Db().Where(
		"user_name = ?", userName,
	).First(&result).Error; err == nil {
		if err := g.Db().Model(result).UpdateColumn(
			"requests_made", gorm.Expr("requests_made + ?", 1),
		).Error; err != nil {
			log.Errorf("user: %s (update number of requests)", err)
		}

		return &result
	} else {
		log.Debugf("user %s not found", txt.Quote(userName))
		return nil
	}
}

// DeleteUserByID removes a given user from the database.
func (u *User) DeleteUserByID() {
	u.DeleteAllUserSummaries()
	g.Db().Delete(&u)
}

// DeleteAllUserSummaries removes all user's summaries from the database
func (u *User) DeleteAllUserSummaries() {
	g.Db().Where("user_id = ?", u.ID).Delete(&Summary{})
}

// AllSummaries returns all summaries of this user
func (u *User) AllSummaries() (result Summaries) {
	g.Db().Where("user_id = ?", u.ID).Order("target_time asc").Find(&result)

	return result
}

// AddSummary will create a new Summary append it to the user entity
func (u *User) AddSummary(
	summary string, videoID string, language string,
) (s Summary, err error) {
	s = Summary{
		ID:       rnd.TLDWID('s'),
		UserID:   u.ID,
		VideoID:  videoID,
		Summary:  summary,
		Score:    0,
		Language: language,
	}

	if err := g.Db().Model(&u).Association("Summaries").Append(&s); err != nil {
		if err := g.Db().Model(u).UpdateColumn(
			"requests_made", gorm.Expr("requests_made + ?", 1),
		).Error; err != nil {
			log.Errorf("user: %s (update number of requests)", err)
		}

		return s, err
	} else {
		return s, nil
	}
}

// String returns an identifier that can be used in logs.
func (u *User) String() string {
	return u.ID
}

// Registered returns true if the user has a user name.
func (u *User) Registered() bool {
	return u.UserName != "" && rnd.IsDCID(u.ID, 'u')
}

// Anonymous returns true if the user is unknown.
func (u *User) Anonymous() bool {
	return !rnd.IsDCID(u.ID, 'u') || u.ID == UnknownUser.ID
}

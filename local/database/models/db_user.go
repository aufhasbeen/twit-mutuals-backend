package models

import "time"

// DBUser represents the data from user destined to enter the database only
type DBUser struct {
	UserID            int64 `gorm:"primaryKey"`
	Token             string
	TokenSecret       string
	TimeTokenInserted string
	Registered        bool
	ScreenName        string

	Mutuals []DBUser `gorm:"many2many:mutual_refer"`
}

// AddMutual adds a single mutual to the user list
func (u DBUser) AddMutual(m ...DBUser) {
	u.Mutuals = append(u.Mutuals, m...)
}

func (u User) RefreshToken(token, secret string) {
	u.Token = token
	u.TokenSecret = secret
	u.TimeTokenInserted = time.Now().Format(time.RFC3339)
}

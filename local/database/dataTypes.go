package database

// User tracks the users that have registered to the application and associates
// them to mutuals and their analysis
type User struct {
	UserID int64 `gorm:"primaryKey"`

	Mutuals []User `gorm:"many2many:mutual_refer"`
}

// MutualStatistics tracks the mutuals of User and tracks the number of interactions
// between From and To
type MutualStatistics struct {
	From int64 `gorm:"primarykey"`
	To   int64 `gorm:"primarykey"`

	Likes    uint32
	Retweets uint32
	Replies  uint32

	Total uint32
}

// Note: MutualStatistics may be better served as a test of the bidirectional
// relationship and can thus be the join table

// TODO: use these structures in the API calls instead

// User methods

// AddMutual adds a single mutual to the user list
func (u *User) AddMutual(m ...User) {
	u.Mutuals = append(u.Mutuals, m...)
}

//

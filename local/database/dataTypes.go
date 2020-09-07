package database

// User tracks the users that have registered to the application and associates
// them to mutuals and their analysis
type User struct {
	UserID int64 `gorm:"primary_key"`

	Mutuals []Mutual
}

// Mutual tracks the mutuals of User and tracks the number of interactions
// and User between mutual
type Mutual struct {
	UserID int64 `gorm:"primary_key;auto_increment:false"`

	Likes    int
	Retweets int
	Replies  int

	Total int
}

// TODO: use these structures in the API calls instead

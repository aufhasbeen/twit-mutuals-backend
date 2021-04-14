package models

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

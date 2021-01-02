package twitterapi

// may make more sense to put this in database or its own package

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/aufheben/mutuals-server/local/database"
)

// ApiUsersToDatabaseUsers converts anaconda.user slices to database.user slices for use by the database
func ApiUsersToDatabaseUsers(apiArray []anaconda.User) []database.User {
	var convertedArray []database.User
	for _, user := range apiArray {
		convertedUser := database.User{}
		convertedUser.UserID = user.Id
		convertedArray = append(convertedArray, convertedUser)
	}

	return convertedArray
}

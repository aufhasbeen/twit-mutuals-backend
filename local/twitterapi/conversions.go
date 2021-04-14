package twitterapi

// may make more sense to put this in database or its own package

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/aufheben/mutuals-server/local/database/models"
)

// AnacondaUsersToDatabaseUsers converts anaconda.user slices to database.user slices for use by the database
func AnacondaUsersToDatabaseUsers(apiArray []anaconda.User) []models.User {
	var convertedArray []models.User
	for _, user := range apiArray {
		convertedUser := models.User{}
		convertedUser.UserID = user.Id
		convertedArray = append(convertedArray, convertedUser)
	}

	return convertedArray
}

// AnacondaIDsToDatabaseUsers converts user ID slices to database.user slices for use by the database
func AnacondaIDsToDatabaseUsers(apiArray []int64) []models.DBUser {
	var convertedArray []models.DBUser
	for _, user := range apiArray {
		convertedUser := models.DBUser{}
		convertedUser.UserID = user
		convertedArray = append(convertedArray, convertedUser)
	}

	return convertedArray
}

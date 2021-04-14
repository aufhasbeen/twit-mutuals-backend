package twitterapi

import (
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

// GetUserByID returns a singular user object or the error for handling
func GetUser(screenName string) (anaconda.User, error) {
	return twitter.GetUsersShow(screenName, *new(url.Values))
}

// GetUsersByID converts a collection of IDs to anaconda user objects
func GetUsersByID(IDs []int64) ([]anaconda.User, error) {
	return twitter.GetUsersLookupByIds(IDs, *new(url.Values))
}

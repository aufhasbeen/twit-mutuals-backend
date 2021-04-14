package twitterapi

import (
	"net/url"
	"sort"

	"github.com/ChimeraCoder/anaconda"
	"github.com/aufheben/mutuals-server/local/conf"
	"github.com/aufheben/mutuals-server/local/database/models"
)

var twitter *anaconda.TwitterApi

// Init initializes the twitter client with the user's oauth credentials
func Init(config *conf.Config) {
	auth := config.Authentication
	twitter = anaconda.NewTwitterApiWithCredentials(
		auth.Developer.Oauth,
		auth.Developer.OauthSecret,
		auth.App.Consumer,
		auth.App.ConsumerSecret)

}

// collectFriends collects friends into a slice from a given list and sorts
// them according to their ID numerically. returns in a channel so it may
// be done asynchronously

// intersectSortedUsers return a third slice that includes only elements that are present in both
func intersectSortedUsers(slice1, slice2 []anaconda.User, comp func(int, int) (bool, bool)) []anaconda.User {
	var finalSlice []anaconda.User

	for i, j := 0, 0; i < len(slice1) && j < len(slice2); {
		eq, less := comp(i, j)

		if less {
			i++
		} else if eq {
			finalSlice = append(finalSlice, slice1[i])
			i++
			j++
		} else {
			j++
		}
	}
	return finalSlice
}

func intersectSortedIDs(slice1, slice2 []int64, comp func(int, int) (bool, bool)) []int64 {
	var finalSlice []int64

	for i, j := 0, 0; i < len(slice1) && j < len(slice2); {
		eq, less := comp(i, j)

		if less {
			i++
		} else if eq {
			finalSlice = append(finalSlice, slice1[i])
			i++
			j++
		} else {
			j++
		}
	}
	return finalSlice
}

// GetFollowers returns only the followers list for a single user in the form of an anaconda.User. unsorted
func GetFollowers(screenName string) []int64 {
	queries := make(url.Values)

	queries.Add("screen_name", screenName)

	followers := twitter.GetFollowersIdsAll(queries)
	return collectFollowersID(followers)
}

// GetMutuals retrieves a list of mutuals from a given username
func GetMutuals(screenName string) []int64 {
	following, followers := getMutualsIDs(screenName)

	mutualsSlice := intersectSortedIDs(followers, following, func(i, j int) (bool, bool) {
		return followers[i] == following[j], followers[i] < following[j]
	})

	return mutualsSlice
}

// GetMutualsUser retrieves a list of mutuals using the anaconda user object
func GetMutualsUser(screenName string) []anaconda.User {
	following, followers := getMutualsLists(screenName)

	mutualsSlice := intersectSortedUsers(followers, following, func(i, j int) (bool, bool) {
		return followers[i].Id == following[j].Id, followers[i].Id < following[j].Id
	})

	return mutualsSlice
}

// func getHomeTimeline(screenName string) ([]anaconda.Tweet, error) {
// 	queries := make(url.Values)

// 	queries.Add("screen_name", screenName)
// 	return twitter.GetHomeTimeline(queries)
// }

//
func getMutualsLists(screenName string) ([]anaconda.User, []anaconda.User) {
	queries := make(url.Values)

	queries.Add("screen_name", screenName)
	queries.Add("include_user_entities", "false")
	friends := twitter.GetFriendsListAll(queries)
	followers := twitter.GetFollowersListAll(queries)

	friendsList := collectFriends(friends)
	follwersList := collectFollowers(followers)
	return friendsList, follwersList
}

func getMutualsIDs(screenName string) ([]int64, []int64) {
	queries := make(url.Values)

	queries.Add("screen_name", screenName)
	friends := twitter.GetFriendsIdsAll(queries)
	followers := twitter.GetFollowersIdsAll(queries)

	friendsList := collectFriendsID(friends)
	follwersList := collectFollowersID(followers)
	return friendsList, follwersList
}

func collectFriends(friendsChannel chan anaconda.FriendsPage) []anaconda.User {
	idList := make([]anaconda.User, 0)
	for page := range friendsChannel {
		idList = append(idList, page.Friends...)
	}

	sort.Slice(idList, func(i, j int) bool {
		return idList[i].Id < idList[j].Id
	})

	return idList
}

func collectFriendsID(friendsChannel chan anaconda.FriendsIdsPage) []int64 {
	idList := make([]int64, 0)
	for page := range friendsChannel {
		idList = append(idList, page.Ids...)
	}

	sort.Slice(idList, func(i, j int) bool {
		return idList[i] < idList[j]
	})

	return idList
}

// collectFollowers collects followers into a list from a given list and sorts
// them according to their ID numerically. returns in a channel so it may be
// asynchronous
func collectFollowers(friendsChannel chan anaconda.FollowersPage) []anaconda.User {
	idList := make([]anaconda.User, 0)
	for page := range friendsChannel {
		idList = append(idList, page.Followers...)
	}

	sort.Slice(idList, func(i, j int) bool {
		return idList[i].Id < idList[j].Id
	})

	return idList
}

func collectFollowersID(followersChannel chan anaconda.FollowersIdsPage) []int64 {
	idList := make([]int64, 0)
	for page := range followersChannel {
		idList = append(idList, page.Ids...)
	}

	sort.Slice(idList, func(i, j int) bool {
		return idList[i] < idList[j]
	})

	return idList
}

// GetUnfollowingMutualsSortedIDS returns a list of unfollowing mutuals sorted using the IDs list instead of anaconda.Users
func GetUnfollowingMutualsSortedIDS(databaseList []models.DBUser, followersList []int64) []models.DBUser {
	// assumes databaseList and followersList are sorted
	unfollowedMutuals := make([]models.DBUser, 0)
	comp := func(i, j int) (bool, bool) {
		return databaseList[i].UserID == followersList[j],
			databaseList[j].UserID < followersList[j]
	}
	databaseListLength := len(databaseList)
	followersListLength := len(followersList)

	for i, j := 0, 0; i < databaseListLength && j < followersListLength; {

		eq, less := comp(i, j)

		if eq {
			i++
			j++
		} else if less {
			unfollowedMutuals = append(unfollowedMutuals, databaseList[i])
			i++
		} else {
			j++
		}
	}
	return unfollowedMutuals
}

// GetUnfollowingMutualsSorted something
func GetUnfollowingMutualsSorted(databaseList []models.User, followersList []anaconda.User) []models.User {
	// assumes databaseList and followersList are sorted
	unfollowedMutuals := make([]models.User, 0)
	comp := func(i, j int) (bool, bool) {
		return databaseList[i].UserID == followersList[j].Id,
			databaseList[j].UserID < followersList[j].Id
	}
	databaseListLength := len(databaseList)
	followersListLength := len(followersList)

	for i, j := 0, 0; i < databaseListLength && j < followersListLength; {

		eq, less := comp(i, j)

		if eq {
			i++
			j++
		} else if less {
			unfollowedMutuals = append(unfollowedMutuals, databaseList[i])
			i++
		} else {
			j++
		}
	}
	return unfollowedMutuals
}

// above get mutuals

// func interactionsUser(comparisonUser anaconda.User, analysesUser anaconda.User) {
// 	// TODO: return a list giving an analyses of user interactions
// 	// (likes ftu, retweets ftu, replies ftu, total ftu, total)
// 	// ftu = from target user

// }

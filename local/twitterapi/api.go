package twitterapi

import (
	"net/url"
	"sort"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/aufheben/mutuals-server/local/database"
)

var twitter *anaconda.TwitterApi
var config *twitterAuthConfig

// Configure sets the app and developer authentication fields
func Configure(configuration twitterAuthConfig) {
	config = &configuration
}

// Init initializes the twitter client with the user's oauth credentials
func Init(oauth, oauthSecret string) {
	twitter = anaconda.NewTwitterApiWithCredentials(
		oauth,
		oauthSecret,
		config.App.Consumer,
		config.App.ConsumerSecret)

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

	mutualsSlice := make([]int64, 0)
	mutualsSlice = intersectSortedIDs(followers, following, func(i, j int) (bool, bool) {
		return followers[i] == following[j], followers[i] < following[i]
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

// GetUnfollowingMutualsSorted something
func GetUnfollowingMutualsSorted(databaseList []database.User, followersList []anaconda.User) []database.User {
	// assumes databaseList and followersList are sorted
	unfollowedMutuals := make([]database.User, 0)
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

//

// below filters for interactions

// GetFilteredRetweets returns a list of twitter users that the authenticated
// user has retweeted.
func GetFilteredRetweets(user anaconda.User) ([]anaconda.Tweet, error) {
	var settings url.Values
	var filteredTweets []anaconda.Tweet

	settings.Add("count", "200")

	for i := 0; i < 800; i += 200 {
		tweets, err := twitter.GetHomeTimeline(settings)

		if err != nil {
			return nil, err
		}

		// removes the user from his own timeline
		filteredTweets = filterTweetsByID(tweets, user.Id)
		lastTweetID := strconv.FormatInt(tweets[len(tweets)-1].Id+1, 10)
		settings.Set("max_id", lastTweetID)
	}

	return filteredTweets, nil
}

// the two get filtered functions are remarkably similar see if you
// can join the two together in some way. not critical

// GetFilteredReplies TODO
func GetFilteredReplies(user anaconda.User) ([]anaconda.Tweet, error) {
	var settings url.Values
	var filteredTweets []anaconda.Tweet

	settings.Add("count", "200")

	for i := 0; i < 800; i += 200 {
		tweets, err := twitter.GetMentionsTimeline(settings)

		if err != nil {
			return nil, err
		}

		// need to change to new function, this one filters by what user
		// posted the reply not to what user the reply was posted to.
		filteredTweets = filterTweetsByID(tweets, user.Id)
		lastTweetID := strconv.FormatInt(tweets[len(tweets)-1].Id+1, 10)
		settings.Set("max_id", lastTweetID)
	}

	return filteredTweets, nil

}

// GetCollectedLikes TODO
func GetCollectedLikes() ([]anaconda.Tweet, error) {
	var settings url.Values
	settings.Add("count", "200")
	var favoriteTweets []anaconda.Tweet

	for i := 0; i < 800; i += 200 {
		favoriteTweets, err := twitter.GetFavorites(settings)

		if err != nil {
			return nil, err
		}

		// need to change to new function, this one filters by what user
		// posted the reply not to what user the reply was posted to.
		lastTweetID := strconv.FormatInt(favoriteTweets[len(favoriteTweets)-1].Id+1, 10)
		settings.Set("max_id", lastTweetID)
	}
	return favoriteTweets, nil
}

// filters tweets by id
func filterTweetsByID(tweets []anaconda.Tweet, ID int64) []anaconda.Tweet {
	filteredTweets := make([]anaconda.Tweet, 0)
	for _, tweet := range tweets {
		if tweet.User.Id != ID {
			filteredTweets = append(filteredTweets, tweet)
		}
	}

	return filteredTweets

}

func interactionsUser(comparisonUser anaconda.User, analysesUser anaconda.User) {
	// TODO: return a list giving an analyses of user interactions
	// (likes ftu, retweets ftu, replies ftu, total ftu, total)
	// ftu = from target user

}

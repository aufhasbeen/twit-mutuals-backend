package twitterapi

import (
	"net/url"
	"sort"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
)

var twitter *anaconda.TwitterApi

// Init initializes the twitter client
func Init() {
	// consumer keys are per user
	// pass them per request
	twitter = anaconda.NewTwitterApiWithCredentials(accessToken, accessTokenSecret, consumerKey, consumerKeySecret)
}

// collectFriends collects friends into a slice from a given list and sorts
// them according to their ID numerically. returns in a channel so it may
// be done asynchronously
func collectFriends(friendsChannel chan anaconda.FriendsPage) chan []anaconda.User {
	list := make(chan []anaconda.User)
	go func() {
		idList := make([]anaconda.User, 0)

		for page := range friendsChannel {
			idList = append(idList, page.Friends...)
		}

		sort.Slice(idList, func(i, j int) bool {
			return idList[i].Id < idList[j].Id
		})

		list <- idList
		close(list)
	}()

	return list
}

// collectFollowers collects followers into a list from a given list and sorts
// them according to their ID numerically. returns in a channel so it may be
// asynchronous
func collectFollowers(friendsChannel chan anaconda.FollowersPage) chan []anaconda.User {
	list := make(chan []anaconda.User)

	go func() {
		idList := make([]anaconda.User, 0)

		for page := range friendsChannel {
			idList = append(idList, page.Followers...)
		}

		sort.Slice(idList, func(i, j int) bool {
			return idList[i].Id < idList[j].Id
		})

		list <- idList
		close(list)
	}()

	return list
}

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
		lastTweetId := strconv.FormatInt(tweets[len(tweets)-1].Id+1, 10)
		settings.Set("max_id", lastTweetId)
	}

	return filteredTweets, nil
}

// the two get filtered functions are remarkably similar see if you
// can join the two together in some way. not critical
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
		lastTweetId := strconv.FormatInt(tweets[len(tweets)-1].Id+1, 10)
		settings.Set("max_id", lastTweetId)
	}

	return filteredTweets, nil

}

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
		lastTweetId := strconv.FormatInt(favoriteTweets[len(favoriteTweets)-1].Id+1, 10)
		settings.Set("max_id", lastTweetId)
	}
	return favoriteTweets, nil
}

// filters tweets by id
func filterTweetsByID(tweets []anaconda.Tweet, Id int64) []anaconda.Tweet {
	filteredTweets := make([]anaconda.Tweet, 0)
	for _, tweet := range tweets {
		if tweet.User.Id != Id {
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

// GetMutuals retrieves a list of mutuals from a given username
func GetMutuals(screenName string) []anaconda.User {
	following, followers := getMutualsLists(screenName)

	mutualsSlice := make([]anaconda.User, 0)
	mutualsSlice = intersectSortedUsers(followers, following, func(i, j int) (bool, bool) {
		return followers[i].Id == following[j].Id, followers[i].Id < following[i].Id
	})

	return mutualsSlice
}

//
func getMutualsLists(screenName string) ([]anaconda.User, []anaconda.User) {
	var queries url.Values

	queries.Add("screen_name", screenName)
	friends := twitter.GetFriendsListAll(queries)
	followers := twitter.GetFollowersListAll(queries)

	friendsList := <-collectFriends(friends)
	follwersList := <-collectFollowers(followers)
	return friendsList, follwersList
}

package twitterapi

import (
	"net/url"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
)

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

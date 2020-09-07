package twitterapi

import (
	"log"
	"sort"

	"github.com/ChimeraCoder/anaconda"
)

type userScorestruct struct {
	user     *User
	total    int
	likes    int
	retweets int
	replies  int
}

type User struct {
	Id         int64
	mutuals    []anaconda.User // needs to be sorted
	OAuthToken string
	topTen     []userScorestruct
}

func (u *User) AnalyzeMutuals() {
	// TODO: take mutuals list and make a list of the top ten mutuals by
	// interaction with the given u
}

func (u *User) makeMutualsMap() map[int64]int {
	mutualsMap := make(map[int64]int)

	for _, mutual := range u.mutuals {
		mutualsMap[mutual.Id] = 0
	}

	return mutualsMap
}

// analyzeMutuals helpers
// much of this function could be moved to reusable helpers
func (u *User) byLikes() {
	type likeIdCount struct {
		id  int64
		val int
	}

	likes, err := GetCollectedLikes()

	// print the error depending on what it may have been
	// most likely is that we did not authenticate
	if err != nil {
		log.Println(err.Error())
		return
	}

	// retrieve a clean map of zeroed out mutuals
	mutualsMap := u.makeMutualsMap()
	mutualsList := make([]likeIdCount, 0)

	// if the tweet was made by a mutual then add one interaction to that
	// mutual
	for _, likeTweet := range likes {
		if _, ok := mutualsMap[likeTweet.User.Id]; ok {
			mutualsMap[likeTweet.User.Id] += 1
		}
	}

	// flatten out the map so that it is a sortable list
	for userId, likeCount := range mutualsMap {
		mutualLikeCount := likeIdCount{userId, likeCount}
		mutualsList = append(mutualsList, mutualLikeCount)
	}

	// sort the list so that u.mutuals and mutualsList have the same user id in
	// the same index
	sort.Slice(mutualsList, func(i int, j int) bool {
		return mutualsList[i].id < mutualsList[j].id
	})

	// add the likes to u.topTen.
	// topTen contains the same users in the same order as mutuals
	// above sounds poorly thought out. top any amount should be a function
	// which retrieves the top amount from u.mutuals.
	// u.mutuals should be changed to be a struct which contains a pointer to
	// anaconda tweets list mutuals and the counts like topten.
	// or keep topten and change name to interactions count and make a function
	// to retrieve the top portions of what was needed.
	for index, likeCount := range mutualsList {
		u.topTen[index].likes = likeCount.val
	}
}

func (u *User) byRetweets() {
}

func (u *User) byReplies() {
}

func (u *User) total() {
}

//

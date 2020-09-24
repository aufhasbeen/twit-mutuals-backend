package routing

import (
	"net/http"
	"strconv"

	"github.com/aufheben/mutuals-server/local/database"
	"github.com/aufheben/mutuals-server/local/twitterapi"
	"github.com/gorilla/mux"
)

// Init initializes the router middleware
func Init() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/unfollows", func(w http.ResponseWriter, r *http.Request) {
		oath := r.Header.Get("oath")
		secret := r.Header.Get("oathSecret")
		user, _ := strconv.ParseInt(r.Header.Get("userID"), 10, 64) // may be variable in url

		twitterapi.InitPerUser(oath, secret)

		// get followers from twitter

		// should also return an error when user is not registered in db
		mutuals := database.FetchMutuals(user)

		// twitterapi.GetUnfollowingMutualsSorted(mutuals, )
	})

	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		oath := r.Header.Get("oath")                                // per user
		secret := r.Header.Get("oathSecret")                        // per user
		user, _ := strconv.ParseInt(r.Header.Get("userID"), 10, 64) // may be variable in url

		// step 1 get user mutuals from twitter

		// step 2  submit to database with oath token and secret for user

	})

	return router
}

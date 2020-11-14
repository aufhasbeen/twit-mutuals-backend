package routing

import (
	"strconv"

	"github.com/aufheben/mutuals-server/local/database"
	"github.com/aufheben/mutuals-server/local/twitterapi"
	"github.com/gin-gonic/gin"
)

var gEngine *gin.Engine

// Init initializes the router middleware
func Init() {
	gEngine = gin.Default()

	gEngine.GET("/unfollows", func(c *gin.Context) {
		oath := c.Request.Header["oath"][0]
		secret := c.Request.Header["oathSecret"][0]
		user, _ := strconv.ParseInt(c.Request.Header["userID"][0], 10, 64) // may be variable in url

		twitterapi.InitPerUser(oath, secret)

		// get followers from twitter

		// should also return an error when user is not registered in db
		mutuals, _ := database.FetchMutuals(user)

		// twitterapi.GetUnfollowingMutualsSorted(mutuals, )
	})

	gEngine.GET("/register", func(c *gin.Context) {
		oath := c.Request.Header["oath"][0]                                // per user
		secret := c.Request.Header["oathSecret"][0]                        // per user
		user, _ := strconv.ParseInt(c.Request.Header["userID"][0], 10, 64) // may be variable in url

		// step 1 get user mutuals from twitter

		// step 2  submit to database with oath token and secret for user

	})

}

// Launch sets the server to running with the default socket
func Launch() {
	gEngine.Run(":8080")
}

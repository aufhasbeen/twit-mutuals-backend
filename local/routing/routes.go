package routing

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/aufheben/mutuals-server/local/database"
	"github.com/aufheben/mutuals-server/local/twitterapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		twitterMutuals := twitterapi.GetMutuals(c.Request.Header["username"][0])

		// should also return an error when user is not registered in db
		dbUser, err := database.FetchUserWithMutuals(user)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(http.StatusNotFound)
			return
		}

		// twitterapi.GetUnfollowingMutualsSorted(mutuals, )
		unfollowers := twitterapi.GetUnfollowingMutualsSorted(dbUser.Mutuals, twitterMutuals)
		c.JSON(http.StatusOK, gin.H{"mutuals": unfollowers})
	})

	gEngine.GET("/register", func(c *gin.Context) {
		oath := c.Request.Header["oath"][0]                                // per user
		secret := c.Request.Header["oathSecret"][0]                        // per user
		user, _ := strconv.ParseInt(c.Request.Header["userID"][0], 10, 64) // may be variable in url

		// step 1 get user mutuals from twitter
		rawMutuals := twitterapi.GetMutuals(c.Request.Header["username"][0])
		dbReadyMutuals := twitterapi.ApiUsersToDatabaseUsers(rawMutuals)

		// step 2  submit to database with oath token and secret for user
		registree := database.User{}
		registree.UserID = user
		registree.Mutuals = dbReadyMutuals
		registree.RefreshToken(oath, secret)

		database.SubmitUser(&registree)
		c.Status(http.StatusOK)
	})

}

// Launch sets the server to running with the default socket
func Launch() {
	gEngine.Run(":8080")
}

package routing

import (
	"errors"
	"net/http"

	"github.com/aufheben/mutuals-server/local/conf"
	"github.com/aufheben/mutuals-server/local/database"
	"github.com/aufheben/mutuals-server/local/database/models"
	"github.com/aufheben/mutuals-server/local/twitterapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var gEngine *gin.Engine

// Init initializes the router middleware
func Init() {

	gEngine = gin.Default()
	api := gEngine.Group("/api/")
	conf.Configure()

	// uncomment when ready for user specific sessions
	// gEngine.Use(func(c *gin.Context) {
	// 	conf.OverwriteUserCredentials(c.Request.Header["oauth"][0], c.Request.Header["oauthSecret"][0])
	// })

	api.GET("/exMutuals", func(c *gin.Context) {
		sessionInfo := sessionHeader{}
		sessionInfo.startSession(c.Request.Header)
		twitterapi.Init(conf.GetConfig())

		screenName := c.Request.URL.Query().Get("screenName")
		if screenName == "" {
			c.JSON(http.StatusExpectationFailed, gin.H{"message": "screenName url paramater not provided"})
			return
		}

		// get followers from twitter
		twitterMutuals := twitterapi.GetMutuals(screenName)

		// should also return an error when user is not registered in db
		dbUser, err := database.FetchUserWithMutuals(sessionInfo.twitUser)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(http.StatusNotFound)
			return
		} else if !dbUser.Registered {
			c.JSON(http.StatusPreconditionFailed, gin.H{"message": "user not registered"})
			return
		}

		// twitterapi.GetUnfollowingMutualsSorted(mutuals, )
		unfollowers := twitterapi.GetUnfollowingMutualsSortedIDS(dbUser.Mutuals, twitterMutuals)
		c.JSON(http.StatusOK, gin.H{"mutuals": unfollowers})
	})

	api.GET("/mutuals", func(c *gin.Context) {
		sessionInfo := sessionHeader{}
		sessionInfo.startSession(c.Request.Header)
		twitterapi.Init(conf.GetConfig())

		screenName := c.Request.URL.Query().Get("screenName")
		if screenName == "" {
			c.JSON(http.StatusExpectationFailed, gin.H{"message": "screenName url paramater not provided"})
			return
		}

		// get followers from twitter
		twitterMutuals := twitterapi.GetMutuals(screenName)

		// twitterapi.GetUnfollowingMutualsSorted(mutuals, )
		c.JSON(http.StatusOK, twitterMutuals)
	})

	api.GET("/profileData", func(c *gin.Context) {
		sessionInfo := sessionHeader{}
		sessionInfo.startSession(c.Request.Header)
		twitterapi.Init(conf.GetConfig())

		screenName := c.Request.URL.Query().Get("screenName")
		println()

		ok, err := twitterapi.GetUser(screenName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, ok)
	})

	api.POST("/profileData", func(c *gin.Context) {
		sessionInfo := sessionHeader{}
		sessionInfo.startSession(c.Request.Header) // currently does nothing
		twitterapi.Init(conf.GetConfig())

		IDs := make([]int64, 0)

		c.BindJSON(&IDs)

		println(len(IDs))
		ok, err := twitterapi.GetUsersByID(IDs)
		if err != nil {
			c.JSON(http.StatusOK, make([]int64, 0))
			return
		}

		c.JSON(http.StatusOK, ok)
	})

	api.GET("/register", func(c *gin.Context) {
		sessionInfo := sessionHeader{}
		sessionInfo.startSession(c.Request.Header)
		twitterapi.Init(conf.GetConfig())

		// step 1 get user mutuals from twitter
		rawMutuals := twitterapi.GetMutuals(c.Query("screenName"))
		dbReadyMutuals := twitterapi.AnacondaIDsToDatabaseUsers(rawMutuals)

		// step 2  submit to database with oath token and secret for user
		registree := models.User{}
		registree.UserID = sessionInfo.twitUser
		registree.Mutuals = dbReadyMutuals
		registree.Registered = true
		// registree.RefreshToken(sessionInfo.oauthPub, sessionInfo.oauthPriv)

		database.SubmitUser(&registree)
		c.Status(http.StatusOK)
	})

	api.GET("/registered", func(c *gin.Context) {
		sessionInfo := sessionHeader{}
		sessionInfo.startSession(c.Request.Header)
		twitterapi.Init(conf.GetConfig())

		screenName := c.Request.URL.Query().Get("screenName")
		if screenName == "" {
			c.JSON(http.StatusExpectationFailed, gin.H{"message": "screenName url paramater not provided"})
			return
		}

		user := database.FetchUser(screenName)
		c.JSON(http.StatusOK, user.Registered)
	})

	api.GET("/unfollow", func(c *gin.Context) {
		sessionInfo := sessionHeader{}
		sessionInfo.startSession(c.Request.Header)

		twitterapi.Init(conf.GetConfig())
	})

}

// Launch sets the server to running with the default socket
func Launch() {
	gEngine.Run(":8080")
}

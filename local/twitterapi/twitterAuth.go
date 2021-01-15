package twitterapi

// Config holds the configuration data for the developer and twitter app
type Config struct {
	Authentication twitterAuthConfig
}

type twitterAuthConfig struct {
	App       twitterAppAuth
	Developer developerAuth
}

type twitterAppAuth struct {
	Consumer       string
	ConsumerSecret string
}

type developerAuth struct {
	Oauth       string
	OauthSecret string
}

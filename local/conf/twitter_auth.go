package conf

import (
	"os"

	"github.com/spf13/viper"
)

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

var conf *Config

// Configure reads in the configuration file containing authentication data per this developer which may include singe user testing credentials
func Configure() {
	conf = new(Config)
	v := viper.New()
	v.SetConfigName(".config.yaml")
	v.AddConfigPath(".")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		print(err.Error())
		os.Exit(1)
	}
	if err := v.Unmarshal(conf); err != nil {
		print(err.Error())
		os.Exit(2)
	}
}

// OverwriteUserCredentials overwrites the default testing credentials filled out in the configuration file
func OverwriteUserCredentials(oauth string, oauthSecret string) {
	conf.Authentication.Developer.Oauth = oauth
	conf.Authentication.Developer.OauthSecret = oauthSecret
}

// GetConfig gets the configuration object conf
func GetConfig() *Config {
	return conf
}

package twitterapi

import (
	"os"
	"testing"

	"github.com/ChimeraCoder/anaconda"
	"github.com/spf13/viper"
)

func setup() {
	var Conf Config

	viper.SetConfigName(".config.yaml")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		print(err.Error())
		os.Exit(1)
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		print(err.Error())
		os.Exit(2)
	}
	authConfig := Conf.authentication

	Configure(authConfig)
	Init(authConfig.Developer.Oauth, authConfig.Developer.OauthSecret)
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
}

func TestSliceReflectionUser(test *testing.T) {
	print("here")
	arr1 := make([]anaconda.User, 6)
	arr2 := make([]anaconda.User, 3)

	for x := range arr1 {
		arr1[x].Id = int64(x)
	}

	for x := range arr2 {
		arr2[x].Id = int64(x + 3)
	}

	result := intersectSortedUsers(arr1, arr2, func(i, j int) (bool, bool) {
		return arr1[i].Id == arr2[j].Id, arr1[i].Id < arr2[j].Id
	})

	if !(result[0].Id == 3 && result[1].Id == 4 && result[2].Id == 5) {
		test.Errorf("Result was %#v", arr2)
	}
}

// func TestGetMutuals(test *testing.T) {
// 	screenName := ""

// 	mutuals := GetMutuals(screenName)

// 	if len(mutuals) > 0 {
// 		print(mutuals[0])
// 	}
// }
